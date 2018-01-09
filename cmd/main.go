package main

import (
	"fmt"
	"github.com/tchaudhry91/podcast-manage-svc/podcastmg"
)

func main() {
	dbStore := podcastmg.NewDBStore(
		"postgres",
		"host=localhost user=test password=password sslmode=disable dbname=podcastmg",
	)
	err := dbStore.Connect()
	if err != nil {
		fmt.Errorf(err.Error())
	}
	fmt.Println("Connected to DB")
	dbStore.DropExistingTables()
	fmt.Println("Dropped Existing Tables")
	dbStore.Migrate()
	fmt.Println("Migration complete")
	user := podcastmg.User{UserEmail: "tc@test.com"}
	podcastItem := podcastmg.PodcastItem{Title: "Episode1"}
	podcastItem2 := podcastmg.PodcastItem{Title: "Episode2"}
	podcast := podcastmg.Podcast{Title: "BeyondID", PodcastItems: []podcastmg.PodcastItem{podcastItem, podcastItem2}, URL: "http://beyond.com/xml"}
	dbStore.CreateUser(&user)
	dbStore.CreatePodcast(&podcast)

	var testUser podcastmg.User
	testUser, _ = dbStore.GetUserFromEmail("tc@test.com")
	fmt.Println(testUser)
	testUser.AddSubscription(podcast)
	dbStore.UpdateUser(&testUser)

	var testUser2 podcastmg.User
	testUser2, _ = dbStore.GetUserFromEmail("tc@test.com")
	fmt.Println(testUser2)
	fmt.Println(testUser2.GetSubscriptions())

}
