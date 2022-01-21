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
	"strings"
)

var appstoreRatingRegex = regexp.MustCompile(`(\d\.\d) • ([\d.]+)(K?) Ratings`)

func GetAppstoreInfo(link string) (ExtraInfo, error) {
	var result AppstoreInfo

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
	if ratingText != "" {
		submatch := appstoreRatingRegex.FindStringSubmatch(ratingText)
		if len(submatch) != 4 {
			return nil, errors.New(fmt.Sprintf("Unparsable rating text %s", ratingText))
		}

		result.Ratings, err = strconv.ParseFloat(submatch[2], 8)
		if err != nil {
			log.Fatalf("Not a number %s", submatch[2])
		}

		if submatch[3] == "K" {
			result.Ratings *= 1000
		}

		result.Score, err = strconv.ParseFloat(submatch[1], 8)
		if err != nil {
			log.Fatalf("Not a number %s", submatch[1])
		}
	}

	reader.Find(".information-list__item__term").EachWithBreak(
		func(_ int, selection *goquery.Selection) bool {
			if dt := selection.Get(0); dt.FirstChild.Data == "Category" {
				a := dt.NextSibling.NextSibling.FirstChild.NextSibling
				if a.Data != "a" {
					panic(fmt.Sprintf("Unexpected node %s", a.Data))
				}

				result.Category = strings.TrimSpace(a.FirstChild.Data)
				return false
			}

			return true
		})

	return result, nil

}

type AppstoreInfo struct {
	Score    float64
	Ratings  float64
	Category string
}

func (a AppstoreInfo) GetLabels() []Label {
	return []Label{
		{"Score", fmt.Sprintf("%.1f", a.Score)},
		{"Ratings", fmt.Sprintf("%.0f", a.Ratings)},
	}
}
