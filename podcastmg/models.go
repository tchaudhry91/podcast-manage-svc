package podcastmg

import (
	"errors"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
	"time"
)

// User is a struct that holds information of a User
type User struct {
	gorm.Model `json:"-"`
	UserEmail  string    `gorm:"not null; unique" json:"user_email"`
	password   string    `gorm:"not null;" json:"-"`
	admin      bool      `json:"-"`
	Podcasts   []Podcast `gorm:"many2many:subscriptions;" json:"podcasts"`
}

// NewUser constructs a User struct with the given email and password
func NewUser(email, password string) (User, error) {
	var user User
	if email == "" || password == "" {
		return user, errors.New("Email or password cannot be empty for user")
	}
	passwordHashBytes, err := bcrypt.GenerateFromPassword([]byte(password), 15)
	if err != nil {
		return user, err
	}
	passwordHash := string(passwordHashBytes)
	return User{
		UserEmail: email,
		admin:     false,
		password:  passwordHash,
	}, nil
}

// Podcast is a struct containing information relevant to a particular podcast
type Podcast struct {
	gorm.Model   `json:"-"`
	PodcastItems []PodcastItem `json:"podcast_items"`
	Title        string        `gorm:"not null" json:"title"`
	Description  string        `json:"description"`
	ImageURL     string        `json:"image_url"`
	URL          string        `gorm:"not null; unique" json:"url"`
}

// NewPodcast constructs a Podcast struct with the given parameters
func NewPodcast(title, description, imageURL, feedURL string, items []PodcastItem) Podcast {
	return Podcast{
		Title:        title,
		Description:  description,
		ImageURL:     imageURL,
		URL:          feedURL,
		PodcastItems: items,
	}
}

// PodcastItem is a struct representing a single item in a given podcast
type PodcastItem struct {
	gorm.Model  `json:"-"`
	PodcastId   uint       `gorm:"index" json:"podcast_id"`
	Title       string     `json:"title"`
	Content     string     `json:"content"`
	Description string     `json:"description"`
	MediaURL    string     `json:"media_url"`
	MediaLength string     `json:"media_length"`
	ImageURL    string     `json:"image_url"`
	Published   *time.Time `json:"published"`
}

// NewPodcastItem constructs a PodcastItem struct with the given values
func NewPodcastItem(title, description, content, mediaURL, imageURL, mediaLength string, published *time.Time) PodcastItem {
	return PodcastItem{
		Title:       title,
		Description: description,
		Content:     content,
		MediaURL:    mediaURL,
		MediaLength: mediaLength,
		ImageURL:    imageURL,
		Published:   published,
	}
}

// GetParentID returns the PodcastID of the podcast that the Item is a part of
func (podcastItem *PodcastItem) GetParentId() uint {
	return podcastItem.PodcastId
}

// GetItems returns the slice of PodcastItem which belongs to this podcast
func (podcast *Podcast) GetItems() []PodcastItem {
	return podcast.PodcastItems
}

// AddSubscriptions adds a podcast to the user's slice of subscribed podcasts
func (user *User) AddSubscription(podcast Podcast) {
	user.Podcasts = append(user.Podcasts, podcast)
}

// GetSubscriptions returns a slice of podcasts that the user is subscribed to
func (user *User) GetSubscriptions() []Podcast {
	return user.Podcasts
}

// GetUserEmail returns the user's email id
func (user *User) GetUserEmail() string {
	return user.UserEmail
}

// ComparePasswords compares the user's hashed password to the given password, returns nil on success
func (user *User) ComparePassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(user.password), []byte(password))
}
