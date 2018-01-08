package main

import (
	"fmt"
	"github.com/tchaudhry91/podcast-manage-svc/podcastmg"
)

func main() {
	dbStore := podcastmg.NewDBStore(
		"postgres",
		"host=localhost user=testuser password=password sslmode=disable dbname=podcastmg",
	)
	err := dbStore.Connect()
	if err != nil {
		fmt.Errorf(err.Error())
	}
	fmt.Println("Connected to DB")
	dbStore.Migrate()
	fmt.Println("Migration complete")
	user := podcastmg.User{UserEmail: "tc@test.com"}
	podcastItem := podcastmg.PodcastItem{PodcastId: "BeyondID", PodcastItemId: "Episode1"}
	podcast := podcastmg.Podcast{PodcastId: "BeyondID", PodcastItems: []podcastmg.PodcastItem{podcastItem}}
	sub := podcastmg.Subscription{PodcastId: "BeyondID", UserEmail: "tc@test.com"}
	dbStore.CreateUser(&user)
	dbStore.CreatePodcast(&podcast)
	dbStore.CreatePodcastItem(&podcastItem)
	dbStore.CreateSubscription(&sub)

	var podcastItems []podcastmg.PodcastItem
	//dbStore.Database.Model(&podcast).Related(&podcastItems)
	dbStore.Database.Model(&podcast).Association("podcast_items").Find(&podcastItems)
	fmt.Println(podcastItems)

	var subs []podcastmg.Subscription
	dbStore.Database.Model(&podcast).Related(&subs)
	fmt.Println(subs)

	var subs2 []podcastmg.Subscription
	dbStore.Database.Model(&user).Related(&subs2)
	fmt.Println(subs2)
}
