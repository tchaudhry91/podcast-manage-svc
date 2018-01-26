package podcastmg

import (
	"errors"
	"github.com/jinzhu/gorm"
	"time"
)

type User struct {
	gorm.Model `json:"-"`
	UserEmail  string    `gorm:"not null; unique" json:"user_email"`
	Admin      bool      `json:"admin"`
	Podcasts   []Podcast `gorm:"many2many:subscriptions;" json:"podcasts"`
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
	gorm.Model   `json:"-"`
	PodcastItems []PodcastItem `json:"podcast_items"`
	Title        string        `gorm:"not null" json:"title"`
	Description  string        `json:"description"`
	ImageURL     string        `json:"image_url"`
	URL          string        `gorm:"not null; unique" json:"url"`
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
