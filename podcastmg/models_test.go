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
