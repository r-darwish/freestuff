package freestuff

import (
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io"
	"log"
	"net/http"
	"regexp"
	"strconv"
)

var regex = regexp.MustCompile(`(\d\.\d) • ([\d.]+)(K?) Ratings`)

func GetAppstoreInfo(link string) (ExtraInfo, error) {
	response, err := http.Get(link)
	if err != nil {
		return nil, err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Printf("Closing repsonse from apple.com: %v", err.Error())
		}
	}(response.Body)

	if !ResponseOk(response.StatusCode) {
		return nil, errors.New(fmt.Sprintf("Status code %s: %d", link, response.StatusCode))
	}

	reader, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		return nil, err
	}

	ratingText := reader.Find(".we-rating-count").Text()
	if ratingText == "" {
		return nil, errors.New(fmt.Sprintf("Could not find the rating section in %s", link))
	}

	submatch := regex.FindStringSubmatch(ratingText)
	if len(submatch) != 4 {
		return nil, errors.New(fmt.Sprintf("Unparsable rating text %s", ratingText))
	}

	ratings, err := strconv.ParseFloat(submatch[2], 8)
	if err != nil {
		log.Fatalf("Not a number %s", submatch[2])
	}

	if submatch[3] == "K" {
		ratings *= 1000
	}

	score, err := strconv.ParseFloat(submatch[1], 8)
	if err != nil {
		log.Fatalf("Not a number %s", submatch[1])
	}

	return &AppstoreInfo{
		Score:   score,
		Ratings: ratings,
	}, nil

}

type AppstoreInfo struct {
	Score   float64
	Ratings float64
}

func (a AppstoreInfo) GetLabels() []Label {
	return []Label{
		{"Score", fmt.Sprintf("%.1f", a.Score)},
		{"Ratings", fmt.Sprintf("%.0f", a.Ratings)},
	}
}
