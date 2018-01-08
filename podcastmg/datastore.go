package podcastmg

import (
	"errors"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type User struct {
	UserEmail string `gorm:"primary_key"`
}

type Podcast struct {
	PodcastId    string        `gorm:"primary_key"`
	PodcastItems []PodcastItem `gorm:"ForeignKey:PodcastId"`
	Title        string
	Description  string
	ImageURL     string
	Tags         string
	URL          string
}

type PodcastItem struct {
	PodcastId     string
	PodcastItemId string `gorm:"primary_key"`
	Title         string
	MediaURL      string
	ImageURL      string
}

type Subscription struct {
	PodcastId string `gorm:"primary_key"`
	UserEmail string `gorm:"primary_key"`
}

func (dbStore *DBStore) Connect() error {
	db, err := gorm.Open(dbStore.dialect, dbStore.connectionString)
	if err != nil {
		return err
	}
	dbStore.Database = db
	return nil
}

func (dbStore *DBStore) Migrate() error {
	if err := dbStore.Database.AutoMigrate(&Podcast{}, &User{}, &PodcastItem{}, &Subscription{}).Error; err != nil {
		return err
	}
	dbStore.Database.Model(&PodcastItem{}).AddForeignKey(
		"podcast_id",
		"podcasts(podcast_id)",
		"RESTRICT",
		"RESTRICT",
	)
	dbStore.Database.Model(&Subscription{}).AddForeignKey(
		"podcast_id",
		"podcasts(podcast_id)",
		"RESTRICT",
		"RESTRICT",
	)
	dbStore.Database.Model(&Subscription{}).AddForeignKey(
		"user_email",
		"users(user_email)",
		"RESTRICT",
		"RESTRICT",
	)
	return nil
}

func (dbStore *DBStore) CreateUser(user *User) error {
	if ok := dbStore.Database.NewRecord(user); !ok {
		errors.New("User Already Exists")
	}
	if err := dbStore.Database.Create(user).Error; err != nil {
		return err
	}
	return nil
}

func (dbStore *DBStore) CreatePodcast(podcast *Podcast) error {
	if err := dbStore.Database.Create(podcast).Error; err != nil {
		return err
	}
	return nil
}

func (dbStore *DBStore) CreatePodcastItem(podcastItem *PodcastItem) error {
	if err := dbStore.Database.Create(podcastItem).Error; err != nil {
		return err
	}
	return nil
}

func (dbStore *DBStore) CreateSubscription(sub *Subscription) error {
	if err := dbStore.Database.Create(sub).Error; err != nil {
		return err
	}
	return nil
}

func NewDBStore(dialect string, connectionString string) *DBStore {
	dbStore := DBStore{
		dialect,
		connectionString,
		nil,
	}
	return &dbStore
}

type DBStore struct {
	dialect          string
	connectionString string
	Database         *gorm.DB
}
