package podcastmg

import (
	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model
	UserEmail string    `gorm:"not null; unique"`
	Podcasts  []Podcast `gorm:"many2many:subscriptions;"`
}

type Podcast struct {
	gorm.Model
	PodcastItems []PodcastItem
	Title        string `gorm:"not null"`
	Description  string
	ImageURL     string
	Tags         string
	URL          string `gorm:"not null; unique"`
}

type PodcastItem struct {
	gorm.Model
	PodcastId uint `gorm:"index"`
	Title     string
	MediaURL  string
	ImageURL  string
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
