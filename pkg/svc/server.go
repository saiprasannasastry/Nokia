package svc

import (
	"context"
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres" //importing postgres

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	albumgpb "github.com/Nokia/proto"
	"github.com/Shopify/sarama"
)

//AlbumServiceServer returns an instance of db and grpc server
type AlbumServiceServer struct {
	db    *gorm.DB
	kafka sarama.SyncProducer
}

//Album maps the fields present in postgres
type Album struct {
	Title        string `json:"title" gorm:"size:255;not null;unique"`
	Albumid      int    `json:"albumid"`
	Id           int    `json:"id" gorm:"primary_key;not null"`
	Url          string `json:"url" gorm:"size:255;not null;unique"`
	Thumbnailurl string `json:"thumbnailurl"`
}

//NewAlbumServer returns an instance of AlbumServiceServer
func NewAlbumServer(db *gorm.DB, kafka sarama.SyncProducer) albumgpb.AlbumServiceServer {
	return &AlbumServiceServer{db: db, kafka: kafka}
}

func (a *AlbumServiceServer) publishKafka(value string) {

	topic := "test"
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(value),
	}

	partition, offset, err := a.kafka.SendMessage(msg)
	if err != nil {
		panic(err)
	}

	log.Infof("Message is stored in topic(%s)/partition(%d)/offset(%d)\n", topic, partition, offset)
}

//CreateAlbum creates a photo with the given id and returns error if same id is present
func (a *AlbumServiceServer) CreateAlbum(ctx context.Context, in *albumgpb.Albumreq) (*albumgpb.CreateAlbumResponse, error) {
	tx := a.db.Begin()

	photo := Album{}
	photo.Title = in.Album.Title
	photo.Albumid = int(in.Album.AlbumId)
	photo.Id = int(in.Album.Id)
	photo.Url = in.Album.Url
	photo.Thumbnailurl = in.Album.ThumbNailUrl

	if err := tx.Table("public.photo").Create(&photo); err.Error != nil {
		st := status.New(codes.AlreadyExists, "ALREADY_EXISTS")
		log.Errorf("%+v already exists", photo)
		tx.Rollback()
		return nil, st.Err()
	}
	tx.Commit()
	log.Infof("succesfully added Photo with id:%v and title:%v into the database", photo.Id, photo.Title)

	message := fmt.Sprintf("method:POST request received with value %+v", photo)
	a.publishKafka(message)

	return &albumgpb.CreateAlbumResponse{Message: "Succesfully added photo into DB"}, nil
}

//GetAlbums returns list of all albums
func (a *AlbumServiceServer) GetAlbums(test *empty.Empty, in albumgpb.AlbumService_GetAlbumsServer) error {

	rows, err := a.db.Table("photo").Select("*").Rows()

	if err != nil {
		return err
	}

	var album Album
	for rows.Next() {
		err := rows.Scan(&album.Id, &album.Albumid, &album.Title, &album.Url, &album.Thumbnailurl)
		if err != nil {
			log.Infof("Error scanning rows")
			return err
		}
		log.Infof("The album sent to client are %+v", album)
		in.Send(&albumgpb.Albumreq{
			Album: &albumgpb.Photo{
				Id:           int64(album.Id),
				AlbumId:      int64(album.Albumid),
				Title:        album.Title,
				Url:          album.Url,
				ThumbNailUrl: album.Thumbnailurl,
			},
		})
	}
	msg := fmt.Sprintf("method:GET request received to send all albums ")
	a.publishKafka(msg)

	return nil
}

//GetAlbum by id
func (a *AlbumServiceServer) GetAlbum(in *albumgpb.GetAlbumreqParams, stream albumgpb.AlbumService_GetAlbumServer) error {

	rows, err := a.db.Table("photo").Select("*").Where("albumid = ?", int(in.AlbumId)).Rows()

	if err != nil {
		log.Errorf("Failed to get Rows %v", err)
		return err
	}

	var album Album
	for rows.Next() {
		err := rows.Scan(&album.Id, &album.Albumid, &album.Title, &album.Url, &album.Thumbnailurl)
		if err != nil {
			//set error code
			msg := "The row does not exists"
			st := status.New(codes.Unavailable, msg)
			log.Errorf(msg)
			return st.Err()
		}

		stream.Send(&albumgpb.Albumreq{
			Album: &albumgpb.Photo{
				Id:           int64(album.Id),
				AlbumId:      int64(album.Albumid),
				Title:        album.Title,
				Url:          album.Url,
				ThumbNailUrl: album.Thumbnailurl,
			},
		})
	}
	log.Infof("The particular album sent is %+v", int(in.AlbumId))
	msg := fmt.Sprintf("method:GET request received to send  album with id %v", int(in.AlbumId))
	a.publishKafka(msg)

	return nil
}

//GetPhoto after selecting the album
func (a *AlbumServiceServer) GetPhoto(ctx context.Context, in *albumgpb.GetphotoReq) (*albumgpb.Photo, error) {

	var album Album
	rows := a.db.Table("photo").Select("*").Where("id = ? AND albumid = ?", int(in.PhotoId), int(in.AlbumId)).Row()

	err := rows.Scan(&album.Id, &album.Albumid, &album.Title, &album.Url, &album.Thumbnailurl)
	if err != nil {
		//set error code
		msg := "The row does not exists"
		st := status.New(codes.Unavailable, msg)
		log.Errorf(msg)
		return nil, st.Err()
	}
	msg := fmt.Sprintf("method:GET request received to send  album with id:%v and photo id:%v", int(in.AlbumId), int(in.PhotoId))
	a.publishKafka(msg)

	return &albumgpb.Photo{Id: int64(album.Id),
			AlbumId:      int64(album.Albumid),
			Title:        album.Title,
			Url:          album.Url,
			ThumbNailUrl: album.Thumbnailurl},
		nil
}

//UpdatePhoto moves the location of photo from one id to another
func (a *AlbumServiceServer) UpdatePhoto(ctx context.Context, in *albumgpb.UpdatePhotoReq) (*empty.Empty, error) {

	var album Album
	rows := a.db.Table("photo").Select("*").Where("title = ? AND albumid = ?", in.OldTitle, int(in.OldAlbumId)).Row()

	err := rows.Scan(&album.Id, &album.Albumid, &album.Title, &album.Url, &album.Thumbnailurl)
	if err != nil {
		//set error code
		msg := "The row does not exists"
		st := status.New(codes.Unavailable, msg)
		log.Errorf(msg)
		return nil, st.Err()
	}
	tx := a.db.Begin()
	if err = tx.Table("public.photo").Delete(&album).Error; err != nil {
		tx.Rollback()
		return nil, err
	}
	album.Albumid = int(in.NewAlbumId)
	album.Title = in.NewTitle
	if err := tx.Table("public.photo").Create(&album); err.Error != nil {
		msg := fmt.Sprintf("value with id:%v , title:%v, url:%v ,thumbnailurl:%v already exists", album.Id, album.Title, album.Url, album.Thumbnailurl)
		st := status.New(codes.AlreadyExists, msg)
		log.Errorf(msg)
		tx.Rollback()
		return nil, st.Err()
	}
	tx.Commit()

	msg := fmt.Sprintf("method:PUT request received ")
	a.publishKafka(msg)

	return new(empty.Empty), nil
}

//DeleteAlbum deletes the full album
func (a *AlbumServiceServer) DeleteAlbum(ctx context.Context, in *albumgpb.DeleteReq) (*empty.Empty, error) {
	var photo Album

	tx := a.db.Begin()
	row := a.db.Table("public.photo").Select("*").Where("id=?", int(in.PhotoId)).Row()
	if row == nil {
		//set error code
		msg := "The row does not exists"
		st := status.New(codes.Unavailable, msg)
		log.Errorf(msg)
		return nil, st.Err()
	}
	if err := tx.Table("public.photo").Where("id=?", int(in.PhotoId)).Delete(&photo).Error; err != nil {
		tx.Rollback()
		log.Errorf("Could not delete because :%v ", err)
		return nil, err
	}
	tx.Commit()

	msg := fmt.Sprintf("method:DELETE request received  to delete photo with id:%v", int(in.PhotoId))
	a.publishKafka(msg)

	return new(empty.Empty), nil
}
