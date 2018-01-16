package podcastmg

import (
	"errors"
	"github.com/jinzhu/gorm"
	"time"
)

type User struct {
	gorm.Model
	UserEmail string `gorm:"not null; unique"`
	Admin     bool
	Podcasts  []Podcast `gorm:"many2many:subscriptions;"`
}

func NewUser(email string, admin bool) (User, error) {
	var user User
	if email == "" {
		return user, errors.New("Email cannot be empty for user")
	}
	return User{
		UserEmail: email,
		Admin:     admin,
	}, nil
}

type Podcast struct {
	gorm.Model
	PodcastItems []PodcastItem
	Title        string `gorm:"not null"`
	Description  string
	ImageURL     string
	URL          string `gorm:"not null; unique"`
}

func NewPodcast(title, description, imageURL, feedURL string, items []PodcastItem) Podcast {
	return Podcast{
		Title:        title,
		Description:  description,
		ImageURL:     imageURL,
		URL:          feedURL,
		PodcastItems: items,
	}
}

type PodcastItem struct {
	gorm.Model
	PodcastId   uint `gorm:"index"`
	Title       string
	Content     string
	Description string
	MediaURL    string
	MediaLength string
	ImageURL    string
	Published   *time.Time
}

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

func (podcastItem *PodcastItem) GetParentId() uint {
	return podcastItem.PodcastId
}

func (podcast *Podcast) GetItems() []PodcastItem {
	return podcast.PodcastItems
}

func (user *User) AddSubscription(podcast Podcast) {
	user.Podcasts = append(user.Podcasts, podcast)
}

func (user *User) GetSubscriptions() []Podcast {
	return user.Podcasts
}

func (user *User) GetUserEmail() string {
	return user.UserEmail
}
