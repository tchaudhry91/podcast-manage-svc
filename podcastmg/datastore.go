package podcastmg

import (
	"errors"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// Store is an interface that defines the methods needed for a podcast-manage service datastore
type Store interface {
	Connect() error
	Close() error
	Migrate() error
	CleanStore() error
	CreateUser(User) error
	GetUserByEmail(string) (User, error)
	UpdateUser(User) error
	DeleteUserByEmail(string) (User, error)
	GetPodcastById(uint) (Podcast, error)
	CreatePodcast(Podcast) error
}

// Connect creates a connection to the database based on the Store's config. This must be called before any other datastore operations
func (dbStore *DBStore) Connect() error {
	db, err := gorm.Open(dbStore.dialect, dbStore.connectionString)
	if err != nil {
		return err
	}
	db.LogMode(false)
	dbStore.Database = db
	return nil
}

// Close ends the connection to the database
func (dbStore *DBStore) Close() error {
	if dbStore.Database == nil {
		return errors.New("Database object is nil")
	}
	dbStore.Database.Close()
	return nil
}

// Migrate creates database tables and constraints based on the models. This does not delete old structures
func (dbStore *DBStore) Migrate() error {
	if err := dbStore.Database.AutoMigrate(&Podcast{}, &User{}, &PodcastItem{}).Error; err != nil {
		return err
	}
	return nil
}

// DropExistingTables removes old tables completely from the database
func (dbStore *DBStore) DropExistingTables() {
	dbStore.Database.DropTableIfExists(&Podcast{}, &User{}, &PodcastItem{}, "subscriptions")
}

// CleanStore clears the database's existing tables
func (dbStore *DBStore) CleanStore() {
	dbStore.DropExistingTables()
}

// CreateUser creates a user in the database, returns err if user exists
func (dbStore *DBStore) CreateUser(user *User) error {
	if err := dbStore.Database.Create(user).Error; err != nil {
		return err
	}
	return nil
}

// GetUserByEmail returns the user object from the database based on UserEmail
func (dbStore *DBStore) GetUserByEmail(userEmail string) (User, error) {
	var user User
	if err := dbStore.Database.Where("user_email = ?", userEmail).Find(&user).Error; err != nil {
		return user, err
	}
	if err := dbStore.Database.Model(&user).Related(&user.Podcasts, "Podcasts").Error; err != nil {
		return user, err
	}
	return user, nil
}

// UpdateUser updates the particular row in the database
func (dbStore *DBStore) UpdateUser(user *User) error {
	if err := dbStore.Database.Save(user).Error; err != nil {
		return err
	}
	return nil
}

// DeleteUser soft-deletes the particular row in the database
func (dbStore *DBStore) DeleteUser(user *User) error {
	if err := dbStore.Database.Delete(user).Error; err != nil {
		return err
	}
	return nil
}

// DeleteUserByEmail deletes a user row in database based on the emailId
func (dbStore *DBStore) DeleteUserByEmail(email string) error {
	user, err := dbStore.GetUserByEmail(email)
	if err != nil {
		return errors.New("User does not exist")
	}
	return dbStore.DeleteUser(&user)
}

// CreatePodcast create a new Podcast row in the database
func (dbStore *DBStore) CreatePodcast(podcast *Podcast) error {
	if err := dbStore.Database.Create(podcast).Error; err != nil {
		return err
	}
	return nil
}

// GetPodcastById returns a podcast from the database with the corresponding ID
func (dbStore *DBStore) GetPodcastById(podcastId uint) (Podcast, error) {
	var podcast Podcast
	if err := dbStore.Database.Where("id = ?", podcastId).Find(&podcast).Error; err != nil {
		return podcast, err
	}
	if err := dbStore.Database.Model(&podcast).Related(&podcast.PodcastItems, "PodcastItems").Error; err != nil {
		return podcast, err
	}
	return podcast, nil
}

// NewDBStore returns a new DBStore with the dialect and connection string set
func NewDBStore(dialect string, connectionString string) *DBStore {
	dbStore := DBStore{
		dialect,
		connectionString,
		nil,
	}
	return &dbStore
}

// DBStore is a SQL based database store
type DBStore struct {
	dialect          string
	connectionString string
	Database         *gorm.DB
}
