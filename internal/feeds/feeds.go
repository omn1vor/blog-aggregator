package feeds

import (
	"encoding/xml"
	"net/http"
)

type rss struct {
	Channel Channel `xml:"channel"`
}

type Channel struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	Items       []struct {
		Title       string `xml:"title"`
		Link        string `xml:"link"`
		Description string `xml:"description"`
	} `xml:"item"`
}

func Fetch(url string) (*Channel, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	rss := rss{}
	decoder := xml.NewDecoder(res.Body)
	err = decoder.Decode(&rss)
	if err != nil {
		return nil, err
	}

	return &rss.Channel, nil
}
