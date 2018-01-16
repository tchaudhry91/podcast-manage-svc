package podcastmg

import (
	"errors"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func (dbStore *DBStore) Connect() error {
	db, err := gorm.Open(dbStore.dialect, dbStore.connectionString)
	if err != nil {
		return err
	}
	db.LogMode(false)
	dbStore.Database = db
	return nil
}

func (dbStore *DBStore) Close() error {
	if dbStore.Database == nil {
		return errors.New("Database object is nil")
	}
	dbStore.Database.Close()
	return nil
}

func (dbStore *DBStore) Migrate() error {
	if err := dbStore.Database.AutoMigrate(&Podcast{}, &User{}, &PodcastItem{}).Error; err != nil {
		return err
	}
	return nil
}

func (dbStore *DBStore) DropExistingTables() {
	dbStore.Database.DropTableIfExists(&Podcast{}, &User{}, &PodcastItem{}, "subscriptions")
}

func (dbStore *DBStore) CreateUser(user *User) error {
	if err := dbStore.Database.FirstOrCreate(user).Error; err != nil {
		return err
	}
	return nil
}

func (dbStore *DBStore) GetUserFromEmail(userEmail string) (User, error) {
	var user User
	if err := dbStore.Database.Where("user_email = ?", userEmail).Find(&user).Error; err != nil {
		return user, err
	}
	if err := dbStore.Database.Model(&user).Related(&user.Podcasts, "Podcasts").Error; err != nil {
		return user, err
	}
	return user, nil
}

func (dbStore *DBStore) UpdateUser(user *User) error {
	if err := dbStore.Database.Save(user).Error; err != nil {
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

func (dbStore *DBStore) GetPodcast(podcastId uint) (Podcast, error) {
	var podcast Podcast
	if err := dbStore.Database.Where("id = ?", podcastId).Find(&podcast).Error; err != nil {
		return podcast, err
	}
	if err := dbStore.Database.Model(&podcast).Related(&podcast.PodcastItems, "PodcastItems").Error; err != nil {
		return podcast, err
	}
	return podcast, nil
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
