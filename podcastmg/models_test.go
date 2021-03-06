package podcastmg

import (
	"testing"
)

var sampleUsers []User
var samplePodcasts []Podcast
var sampleItems []PodcastItem

func init() {
	samplePodcasts = []Podcast{
		{Title: "Beyond!", URL: "ign.beyond.com/xml", PodcastItems: []PodcastItem{}},
		{Title: "Game Scoop!", URL: "gamescoop.ign.com/xml", PodcastItems: []PodcastItem{{Title: "Episode1"}, {Title: "Episode2"}}},
		{Title: "Unlocked!", URL: "unlocked.ign.com/xml", PodcastItems: []PodcastItem{{Title: "Episode1"}}},
	}

	sampleItems = []PodcastItem{
		{Title: "Episode1", PodcastID: 1},
		{Title: "Episode2", PodcastID: 2},
		{Title: "Episode1"},
	}

	sampleUsers = []User{
		{UserEmail: "a@test.com", admin: false, Password: "123123"},
		{UserEmail: "b@test.com", admin: false, Podcasts: []Podcast{samplePodcasts[1], samplePodcasts[2]}, Password: "test"},
		{UserEmail: "", admin: false, Password: "nothing"},
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

func TestUserpasswordHash(t *testing.T) {
	type userPassTestCase struct {
		emailID       string
		password      string
		errorExpected bool
	}
	testCases := []userPassTestCase{
		{"user1@test.com", "test_pass", false},
		{"user2@test.com", "", true},
		{"user3@wwe.com", "simple", false},
	}
	for _, testCase := range testCases {
		user, err := NewUser(testCase.emailID, testCase.password)
		if err != nil && !testCase.errorExpected {
			t.Errorf("Error Creating user:%v", err)
		}
		if err = user.ComparePassword(testCase.password); err != nil && !testCase.errorExpected {
			t.Errorf("password match failed for %s:%v", testCase.password, err)
		}
	}
}

func TestUserSubscriptions(t *testing.T) {
	type userSubscriptionTestCase struct {
		name   string
		user   User
		append []Podcast
		remove []Podcast
		want   []Podcast
	}

	testCases := []userSubscriptionTestCase{
		{"Empty Append", sampleUsers[1], nil, nil, []Podcast{samplePodcasts[1], samplePodcasts[2]}},
		{"Empty User", sampleUsers[3], nil, nil, nil},
		{"Single Append on Empty User", sampleUsers[2], []Podcast{samplePodcasts[1]}, nil, []Podcast{samplePodcasts[1]}},
		{"Single Append on Empty Subs", sampleUsers[0], []Podcast{samplePodcasts[1]}, nil, []Podcast{samplePodcasts[1]}},
		{"Multi Append on Empty Subs", sampleUsers[0], []Podcast{samplePodcasts[2], samplePodcasts[0]}, nil, []Podcast{samplePodcasts[2], samplePodcasts[0]}},
		{"Remove Sub", sampleUsers[1], nil, []Podcast{samplePodcasts[2]}, []Podcast{samplePodcasts[1]}},
	}

	for _, testCase := range testCases {
		for _, append := range testCase.append {
			testCase.user.AddSubscription(append)
		}
		for _, remove := range testCase.remove {
			testCase.user.RemoveSubscription(remove)
		}
		if have := testCase.user.GetSubscriptions(); !compareSubscriptions(have, testCase.want) {
			t.Errorf("TestCase:%v\tWant:%v\tHave:%v", testCase.name, testCase.want, have)
		}
	}
}

func TestPodcastItems(t *testing.T) {
	type podcastItemTestCase struct {
		item PodcastItem
		want uint
	}

	testCases := []podcastItemTestCase{
		{sampleItems[0], 1},
		{sampleItems[1], 2},
		{sampleItems[2], 0},
	}

	for _, testCase := range testCases {
		if have := testCase.item.GetParentID(); have != testCase.want {
			t.Errorf("Get Podcast Parent ID\thave:%v\twant:%v", have, testCase.want)
		}
	}
}

func TestPodcast(t *testing.T) {
	type podcastGetItemsTestCase struct {
		name    string
		podcast Podcast
		want    []PodcastItem
	}
	testCases := []podcastGetItemsTestCase{
		{"No Episodes", samplePodcasts[0], []PodcastItem{}},
		{"2 Episodes", samplePodcasts[1], []PodcastItem{{Title: "Episode1"}, {Title: "Episode2"}}},
	}

	for _, testCase := range testCases {
		if have := testCase.podcast.GetItems(); len(have) != len(testCase.want) {
			t.Errorf("%s\tHave:%v\tWant:%v", testCase.name, have, testCase.want)
		}
	}
}

func TestPodcastUpdation(t *testing.T) {
	type PodcastTest struct {
		feedURL string
	}
	testCases := []PodcastTest{
		{feedURL: "http://feeds.ign.com/ignfeeds/podcasts/beyond?format=xml"},
	}

	for _, tc := range testCases {
		pc, _ := BuildPodcastFromURL(tc.feedURL)
		lenOld := len(pc.PodcastItems)
		pc.PodcastItems = pc.PodcastItems[:5]
		pc.Update()
		lenNew := len(pc.PodcastItems)
		if lenOld != lenNew {
			t.Errorf("Old Length: %d, New Length: %d", lenOld, lenNew)
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
