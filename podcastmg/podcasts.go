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

// BuildPodcastFromURL returns a populated podcast struct from the feedURL
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

// GetNewItems returns a list of PodcastItems that are present in a new slice
func GetNewItems(feedURL string, old []PodcastItem) (update []PodcastItem, err error) {
	feed, err := parseFeed(feedURL)
	if err != nil {
		return
	}
	new := buildItemsFromFeedItems(feed.Items)

	for i := len(new) - 1; i >= 0; i-- {
		if checkItemIndex(new[i], old) < 0 {
			update = append(update, new[i])
		} else {
			break
		}
	}
	return
}

// checkItemExistence returns the index of an item if it is present in the slice, -1 if not
func checkItemIndex(item PodcastItem, items []PodcastItem) int {
	for i, itm := range items {
		if item.Title == itm.Title && itm.MediaURL == itm.MediaURL {
			return i
		}
	}
	return -1
}
