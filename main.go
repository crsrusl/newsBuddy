package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	db "newsBuddy/pkg"
	"os"
)

type feeds struct {
	Feed []feed
}

type feed struct {
	Channel channel `xml:"channel"`
}

type channel struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	Item        []item `xml:"item"`
}

type item struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Pubdate     string `xml:"pubdate"`
	Comments    string `xml:"comments"`
	Description string `xml:"description"`
}

func main() {
	rss := make(chan feed)

	for _, url := range db.Feeds() {
		go getFeeds(url, rss)
	}

	for range rss {
		newFeed := <-rss
		item := getNews(newFeed.Channel)

		fmt.Println(item)
	}

	os.Exit(1)
}

func getNews(f channel) []item {
	return f.Item
}

func getFeeds(url string, rss chan<- feed) {
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	f, err := unmarshalXML(body)
	if err != nil {
		log.Fatal(err)
	} else {
		rss <- f
	}
}

func unmarshalXML(body []byte) (feed, error) {
	f := &feed{}
	err := xml.Unmarshal([]byte(body), &f)
	if err != nil {
		return *f, err
	}

	return *f, nil
}
