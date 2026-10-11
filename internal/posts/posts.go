package posts

import (
	"fmt"
	"os"
	"sort"
	"time"

	"github.com/mmcdole/gofeed"
)

type Post struct {
	Title   string
	Link    string
	PubDate time.Time
}

func GetPosts() ([]Post, error) {
	url := os.Getenv("RSS_URL")
	if url == "" {
		return nil, fmt.Errorf("SHOW_LATEST_POSTS=true but RSS_URL is not set")
	}
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(url)
	if err != nil {
		return nil, err
	}

	var posts []Post
	for _, item := range feed.Items {
		posts = append(posts, Post{
			Title:   item.Title,
			Link:    item.Link,
			PubDate: pubDate(item),
		})
	}

	// Sort posts by publication date (newest first)
	sort.Slice(posts, func(i, j int) bool {
		return posts[i].PubDate.After(posts[j].PubDate)
	})

	// Return only the first 5 posts
	if len(posts) > 5 {
		posts = posts[:5]
	}

	return posts, nil
}

func pubDate(it *gofeed.Item) time.Time {
	if it.PublishedParsed != nil {
		return *it.PublishedParsed
	}
	return time.Time{}
}
