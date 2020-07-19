package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"sync"

	"github.com/Nokia/pkg/svc"
	albumgpb "github.com/Nokia/proto"
	"github.com/Shopify/sarama"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres" // using postgres sql
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func main() {

	//server.Initialize(os.Getenv("DB_DRIVER"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_PORT"), os.Getenv("DB_HOST"), os.Getenv("DB_NAME"))
	prosgretConname := fmt.Sprintf("dbname=%v password=%v port=%v user=%v host=%v sslmode=disable", os.Getenv("DB_NAME"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"), os.Getenv("DB_HOST"))
	log.Infof("conname is %s", prosgretConname)
	db, err := gorm.Open("postgres", prosgretConname)
	if err != nil {
		log.Errorf("the error is :%v", err)
		panic("Failed to connect to database!")
	}
	defer db.Close()
	go startGRPC(db)

	go startHTTP()

	// Block forever
	var wg sync.WaitGroup
	wg.Add(1)
	wg.Wait()

}

func startGRPC(db *gorm.DB) {
	lis, err := net.Listen("tcp", "0.0.0.0:5566")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()

	producer := setKafka()
	srv := svc.NewAlbumServer(db, producer)
	albumgpb.RegisterAlbumServiceServer(grpcServer, srv)

	log.Println("gRPC server ready...")
	grpcServer.Serve(lis)
}
func startHTTP() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Connect to the GRPC server
	conn, err := grpc.Dial("0.0.0.0:5566", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()

	// Register grpc-gateway
	rmux := runtime.NewServeMux()
	client := albumgpb.NewAlbumServiceClient(conn)
	err = albumgpb.RegisterAlbumServiceHandlerClient(ctx, rmux, client)
	if err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()
	mux.Handle("/", rmux)

	mux.HandleFunc("/swagger.json", serveSwagger)
	fs := http.FileServer(http.Dir("www/swagger-ui"))
	mux.Handle("/swagger-ui/", http.StripPrefix("/swagger-ui", fs))
	log.Println("REST server ready...")

	err = http.ListenAndServe("0.0.0.0:8080", logRequest(mux))
	if err != nil {
		log.Fatal(err)
	}
}

func setKafka() sarama.SyncProducer {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	config.Producer.Return.Successes = true

	brokers := []string{"kafka:9092"}
	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		// Should not reach here
		panic(err)
	}
	return producer

}
func logRequest(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s\n", r.RemoteAddr, r.Method, r.URL)
		handler.ServeHTTP(w, r)
	})
}

func serveSwagger(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "www/swagger.json")
}
