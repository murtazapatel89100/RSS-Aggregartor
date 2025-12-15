package rss

import (
	"encoding/xml"
	"io"
	"net/http"
	"time"
)

type RSSFEED struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Language    string    `xml:"language"`
		Items       []RSSITEM `xml:"item"`
	} `xml:"channel"`
}

type RSSITEM struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func UrlToFeed(url string) (*RSSFEED, error) {
	httpClient := http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := httpClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	rssfeed := &RSSFEED{}
	err = xml.Unmarshal(data, rssfeed)
	if err != nil {
		return nil, err
	}

	return rssfeed, nil
}
