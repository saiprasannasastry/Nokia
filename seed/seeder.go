package seed

import (
	"github.com/Nokia/pkg/svc"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres" // using postgres sql
	log "github.com/sirupsen/logrus"
)

var photos = []svc.Album{
	svc.Album{
		Id:           1,
		Albumid:      1,
		Title:        "This is title 1",
		Url:          "This is Url 1 ",
		Thumbnailurl: "This is thumbnail Url1",
	},
	svc.Album{
		Id:           2,
		Albumid:      1,
		Title:        "This is title 2",
		Url:          "This is Url 2 ",
		Thumbnailurl: "This is thumbnail Url2",
	},
}

func Load(db *gorm.DB) {
	err := db.Debug().DropTableIfExists(&svc.Album{}).Error
	if err != nil {
		log.Fatalf("cannot drop table: %v", err)
	}
	err = db.Debug().AutoMigrate(&svc.Album{}).Error
	if err != nil {
		log.Fatalf("cannot migrate table: %v", err)
	}
	log.Infof("Adding dummy values into db")
	for i, _ := range photos {
		err = db.Debug().Model(&svc.Album{}).Create(&photos[i]).Error
		if err != nil {
			log.Fatalf("cannot seed users table: %v", err)
		}
	}
}
