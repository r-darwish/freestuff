package freestuff

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var priceRegex = regexp.MustCompile(`[$€]?(\d+[.,]\d+)[$€]?`)

func ResponseOk(statusCode int) bool {
	return statusCode >= 200 && statusCode < 300
}

func Price(title string) (float64, error) {
	submatch := priceRegex.FindStringSubmatch(title)
	if len(submatch) < 2 {
		return 0, errors.New(fmt.Sprintf("No price in title: %s", title))
	}

	price, err := strconv.ParseFloat(strings.Replace(submatch[1], ",", ".", 1), 8)
	if err != nil {
		return 0, err
	}

	return price, nil
}
