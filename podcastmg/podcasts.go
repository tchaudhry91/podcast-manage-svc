package podcastmg

import (
	"github.com/mmcdole/gofeed"
)

func parseFeed(xmlURL string) (*gofeed.Feed, error) {
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(xmlURL)
	if err != nil {
		return nil, err
	}
	return feed, nil
}

func buildItemsFromFeedItems(feedItems []*gofeed.Item) []PodcastItem {
	var podcastItems []PodcastItem
	for _, item := range feedItems {
		var mediaURL, mediaLength string
		if len(item.Enclosures) < 1 {
			mediaURL = ""
			mediaLength = "0"
		} else {
			mediaURL = item.Enclosures[0].URL
			mediaLength = item.Enclosures[0].Length
		}
		podcastItem := NewPodcastItem(item.Title, item.Description, item.Content, mediaURL, item.Image.URL, mediaLength, item.PublishedParsed)
		podcastItems = append(podcastItems, podcastItem)
	}
	return podcastItems
}

// BuildPodcastFromUrl returns a populated podcast struct from the feedURL
func BuildPodcastFromURL(feedURL string) (Podcast, error) {
	var pc Podcast
	feed, err := parseFeed(feedURL)
	if err != nil {
		return pc, err
	}
	podcastItems := buildItemsFromFeedItems(feed.Items)
	pc = NewPodcast(feed.Title, feed.Description, feed.Image.URL, feedURL, podcastItems)
	return pc, nil
}
