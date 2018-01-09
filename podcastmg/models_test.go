package podcastmg

import (
	"testing"
)

var sampleUsers []User
var samplePodcasts []Podcast

func init() {
	samplePodcasts = []Podcast{
		{Title: "Beyond!", URL: "ign.beyond.com/xml", PodcastItems: []PodcastItem{}},
		{Title: "Game Scoop!", URL: "gamescoop.ign.com/xml", PodcastItems: []PodcastItem{{Title: "Episode1"}, {Title: "Episode2"}}},
		{Title: "Unlocked!", URL: "unlocked.ign.com/xml", PodcastItems: []PodcastItem{{Title: "Episode1"}}},
	}

	sampleUsers = []User{
		{UserEmail: "a@test.com"},
		{UserEmail: "b@test.com", Podcasts: []Podcast{samplePodcasts[1], samplePodcasts[2]}},
		{UserEmail: ""},
		{},
	}

}

func TestUserEmail(t *testing.T) {
	type userEmailTestCase struct {
		user User
		want string
	}
	testCases := []userEmailTestCase{
		{sampleUsers[0], "a@test.com"},
		{sampleUsers[1], "b@test.com"},
		{sampleUsers[2], ""},
		{sampleUsers[3], ""},
	}

	for _, testCase := range testCases {
		if have := testCase.user.GetUserEmail(); have != testCase.want {
			t.Errorf("Want: %v\t Have:%v\n", testCase.want, have)
		}
	}
}

func TestUserSubscriptions(t *testing.T) {
	type userSubscriptionTestCase struct {
		name   string
		user   User
		append []Podcast
		want   []Podcast
	}

	testCases := []userSubscriptionTestCase{
		{"Empty Append", sampleUsers[1], nil, []Podcast{samplePodcasts[1], samplePodcasts[2]}},
		{"Empty User", sampleUsers[3], nil, nil},
		{"Single Append on Empty User", sampleUsers[2], []Podcast{samplePodcasts[1]}, []Podcast{samplePodcasts[1]}},
		{"Single Append on Empty Subs", sampleUsers[0], []Podcast{samplePodcasts[1]}, []Podcast{samplePodcasts[1]}},
		{"Multi Append on Empty Subs", sampleUsers[0], []Podcast{samplePodcasts[2], samplePodcasts[0]}, []Podcast{samplePodcasts[2], samplePodcasts[0]}},
	}

	for _, testCase := range testCases {
		for _, append := range testCase.append {
			testCase.user.AddSubscription(append)
		}
		if have := testCase.user.GetSubscriptions(); !compareSubscriptions(have, testCase.want) {
			t.Errorf("TestCase:%v\tWant:%v\tHave:%v", testCase.name, testCase.want, have)
		}
	}
}

func compareSubscriptions(sl1 []Podcast, sl2 []Podcast) bool {
	if sl1 == nil && sl2 == nil {
		return true
	}
	if len(sl1) != len(sl2) {
		return false
	}
	for i, v := range sl1 {
		if v.URL != sl2[i].URL {
			return false
		}
	}
	return true
}
