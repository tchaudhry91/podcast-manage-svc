package podcastmg

import (
	"testing"
)

func TestPodcastBuilder(t *testing.T) {
	type PodcastTest struct {
		url          string
		wantTitle    string
		wantMinItems int
		err          bool
	}
	testCases := []PodcastTest{
		{"http://feeds.ign.com/ignfeeds/podcasts/beyond?format=xml", "Podcast Beyond", 200, false},
		{"http://www.buzzsprout.com/3195.rss", "The Cloudcast (.net) - Weekly Cloud Computing Podcast", 300, false},
		{"WRONGURL", "", 0, true},
	}
	for _, testCase := range testCases {
		pc, err := BuildPodcastFromURL(testCase.url)
		if err != nil {
			if testCase.err {
				continue
			}
			t.Errorf("Errored:%v", err)
		} else {
			if testCase.err {
				t.Errorf("Should have errorred but did not:%s", testCase.url)
			}
		}
		haveTitle := pc.Title
		haveMinItems := len(pc.PodcastItems)
		if haveTitle != testCase.wantTitle || haveMinItems < testCase.wantMinItems {
			t.Errorf("TitleWant: %s, TitleHave:%s\n ItemsMinWant:%d, ItemsMinWant:%d",
				testCase.wantTitle, haveTitle, testCase.wantMinItems, haveMinItems)
		}
	}
}
