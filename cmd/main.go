package main

import (
	"fmt"
	"github.com/tchaudhry91/podcast-manage-svc/podcastmg"
	"os"
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
	dbStore.CreateUser(&user)

	podcast, err := podcastmg.BuildPodcastFromURL(os.Args[1])
	if err != nil {
		panic(err)
	}
	dbStore.CreatePodcast(&podcast)

	var testUser podcastmg.User
	testUser, _ = dbStore.GetUserFromEmail("tc@test.com")
	fmt.Println(testUser)
	testUser.AddSubscription(podcast)
	dbStore.UpdateUser(&testUser)

	var testUser2 podcastmg.User
	testUser2, _ = dbStore.GetUserFromEmail("tc@test.com")
	subs := testUser2.GetSubscriptions()
	var pc podcastmg.Podcast
	pc, err = dbStore.GetPodcast(subs[0].ID)
	if err != nil {
		panic(err)
	}
	fmt.Println(pc)
}
