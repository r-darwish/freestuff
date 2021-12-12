package freestuff

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"

	"github.com/PuerkitoBio/goquery"
)

type Subreddit string

type RedditLink struct {
	Title      string
	Link       string
	RedditLink string
	Image      string
}

func GetLinksFromSubreddit(subreddit Subreddit) ([]RedditLink, error) {
	client := http.Client{}
	request, err := http.NewRequest("GET", fmt.Sprintf("https://old.reddit.com/r/%s", subreddit), nil)
	if err != nil {
		return nil, err
	}

	request.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.97 Safari/537.36")
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			log.Printf("Error closing the http request: %v", err)
		}
	}(response.Body)

	if response.StatusCode != 200 {
		if response.StatusCode == 429 {
			return nil, errors.New(fmt.Sprintf("Too many requests to reddit. Need to wait %s seconds", response.Header.Get("Retry-After")))
		}

		return nil, errors.New(fmt.Sprintf("Reddit returned %d", response.StatusCode))
	}

	reader, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		return nil, err
	}

	selection := reader.Find(".link")
	result := make([]RedditLink, selection.Length())
	selection.Each(func(i int, selection *goquery.Selection) {
		log.Printf("%v", selection.Text())

		_, isAd := selection.Attr("data-adserver-impression-id")
		if isAd {
			return
		}

		link := RedditLink{}

		attr, _ := selection.Attr("data-url")
		link.Link = attr

		attr, _ = selection.Attr("data-permalink")
		link.RedditLink = "https://reddit.com" + attr

		link.Image = getImage(selection)
		link.Title = selection.Find(".title").Get(1).FirstChild.Data

		relevant_platform, err := regexp.MatchString(`^\[.*(?:Mac|Windows).*\]`, link.Title)

		if err != nil {
			log.Printf("Failed matching on platform: %s", link.Title)
			relevant_platform = true // Ignore match in case of error
		}

		if !relevant_platform {
			log.Printf("Skipping irrelevant platform: %s", link.Title)
			return
		}

		result = append(result, link)
	})

	return result, nil
}

func getImage(selection *goquery.Selection) string {
	thumbmails := selection.Find(".thumbnail")
	if thumbmails.Length() == 0 {
		return ""
	}

	attr, _ := thumbmails.Find("img").Attr("src")
	if attr == "" {
		return ""
	}

	return "https:" + attr
}
