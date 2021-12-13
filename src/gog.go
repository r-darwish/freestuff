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

var gogRatingRegex = regexp.MustCompile(`(\d\.\d)/5`)

func GetGogInfo(link string) (ExtraInfo, error) {
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

	ratingText := reader.Find(".rating").Text()
	if ratingText == "" {
		return nil, errors.New(fmt.Sprintf("Could not find the rating section in %s", link))
	}

	submatch := gogRatingRegex.FindStringSubmatch(ratingText)
	if len(submatch) != 2 {
		return nil, errors.New(fmt.Sprintf("Unparsable rating text %s", ratingText))
	}

	score, err := strconv.ParseFloat(submatch[1], 8)
	if err != nil {
		log.Fatalf("Not a number %s", submatch[1])
	}

	return &GogInfo{
		Score: score,
	}, nil

}

type GogInfo struct {
	Score float64
}

func (g GogInfo) GetLabels() []Label {
	return []Label{
		{"Score", fmt.Sprintf("%.1f", g.Score)},
	}
}
