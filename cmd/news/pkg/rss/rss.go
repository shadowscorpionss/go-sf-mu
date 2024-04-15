package rss

import (
	"encoding/json"
	"encoding/xml"
	"log"
	"net/http"
	"os"
	"sf-mu/pkg/models/news"

	"time"

	strip "github.com/grokify/html-strip-tags-go"
)

type Item struct {
	Title   string `xml:"title"`
	Link    string `xml:"link"`
	Content string `xml:"description"`
	PubDate string `xml:"pubDate"`
}

type Channel struct {
	Items []Item `xml:"channel>item"`
}

type config struct {
	Rss           []string `json:"rss"`
	RequestPeriod int      `json:"request_period"`
}

// GoNews reading rss
func GoNews(configURL string, chanPosts chan<- []news.Post, chanErrs chan<- error) error {
	
	file, err := os.Open(configURL)
	if err != nil {
		return err
	}
	var conf config
	err = json.NewDecoder(file).Decode(&conf)
	if err != nil {
		return err
	}

	log.Println("searching rss-channels")
	for i, r := range conf.Rss {
		go func(r string, i int, chanPosts chan<- []news.Post, chanErrs chan<- error) {
			for {
				log.Println("started  goroutine", i, "by link", r)
				p, err := GetRss(r)
				if err != nil {
					chanErrs <- err
					time.Sleep(time.Second * 10)
					continue
				}
				chanPosts <- p
				log.Println("insert posts from goroutine", i, "by link", r)
				log.Println("goroutine ", i, " is wating for next start")
				time.Sleep(time.Duration(conf.RequestPeriod) * time.Second * 15)
			}
		}(r, i, chanPosts, chanErrs)
	}
	return nil
}

func GetRss(url string) ([]news.Post, error) {
	var c Channel
	//rss request
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	err = xml.NewDecoder(res.Body).Decode(&c)
	if err != nil {
		return nil, err
	}

	//convert rss to post
	var posts []news.Post
	for _, i := range c.Items {
		var p news.Post
		p.Title = i.Title
		p.Content = i.Content
		p.Content = strip.StripTags(p.Content)
		p.Link = i.Link

		t, err := time.Parse(time.RFC1123, i.PubDate)
		if err != nil {
			t, err = time.Parse(time.RFC1123Z, i.PubDate)
		}
		if err != nil {
			t, err = time.Parse("Mon, _2 Jan 2006 15:04:05 -0700", i.PubDate)
		}
		if err == nil {
			p.PubTime = t.Unix()
		}

		posts = append(posts, p)
	}
	return posts, nil
}
