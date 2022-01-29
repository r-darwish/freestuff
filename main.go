package main

import (
	"fmt"
	"github.com/gtuk/discordwebhook"
	"github.com/r-darwish/freestuff/src"
	"log"
	"os"
	"strings"
	"time"
)

func main() {
	os.Exit(run())
}

var author = fmt.Sprint("Reddit Comments")
var cache freestuff.RedisCache

func run() int {
	exitCode := 0

	links, err := freestuff.GetLinksFromSubreddit("apphookup")
	if err != nil {
		log.Fatalf("Fatal error: %v", err.Error())
	}

	cache = freestuff.NewRedisCache()

	for _, link := range links {
		if link.Title == "" {
			continue
		}

		err := handleLink(link)
		if err != nil {
			exitCode = 1
			log.Printf("Error in %v: %+v", link, err.Error())
		}
	}

	return exitCode
}

func handleLink(link freestuff.RedditLink) error {
	known, err := cache.IsKnown(link.Link)
	if err != nil {
		return err
	}

	if known {
		return nil
	}

	defer func(cache *freestuff.RedisCache, title string) {
		err := cache.SetKnown(title)
		if err != nil {
			log.Printf("Error storing a link in the cache: %v", err)
		}
	}(&cache, link.Link)

	extraInfo, err, publish := shouldPublish(link)
	if publish {
		return nil
	}

	var fields []discordwebhook.Field
	t := true
	if extraInfo != nil {
		for _, label := range extraInfo.GetLabels() {
			label := label
			fields = append(fields, discordwebhook.Field{
				Name:   &label.Key,
				Value:  &label.Value,
				Inline: &t,
			})
		}

	}

	embed := discordwebhook.Embed{
		Title:     &link.Title,
		Url:       &link.Link,
		Thumbnail: &discordwebhook.Thumbnail{Url: &link.Image},
		Author: &discordwebhook.Author{
			Name: &author,
			Url:  &link.RedditLink,
		},
		Fields: &fields,
	}

	message := discordwebhook.Message{
		Embeds: &[]discordwebhook.Embed{embed},
	}

	time.Sleep(time.Second * 5)

	err = discordwebhook.SendMessage(freestuff.Config.Webhook, message)
	if err != nil && err.Error() != "" {
		return fmt.Errorf("sending a Discord message: %w", err)
	}

	return nil
}

func shouldPublish(link freestuff.RedditLink) (freestuff.ExtraInfo, error, bool) {
	price, err := freestuff.Price(link.Title)
	if err != nil {
		return nil, err, false
	}

	if !isFree(link.Title) || price < freestuff.Config.PriceThreshold {
		return nil, nil, false
	}

	extraInfo, err := freestuff.GetExtraInfo(link.Link)
	if err != nil {
		log.Printf("Error getting extra info for %s: %v", link.Link, err.Error())
		return nil, nil, false
	}

	if appStoreInfo, ok := extraInfo.(freestuff.AppstoreInfo); ok && appStoreInfo.Category != "Games" {
		log.Printf("Skipping category %s for %s", appStoreInfo.Category, link)
		return nil, nil, false
	}
	return extraInfo, nil, true
}

func isFree(title string) bool {
	return strings.Contains(strings.ToLower(title), "free") || strings.Contains(title, "0.00")
}
