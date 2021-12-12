package main

import (
	"fmt"
	"github.com/gtuk/discordwebhook"
	"github.com/r-darwish/freestuff/freestuff"
	"log"
	"os"
	"strings"
	"time"
)

func main() {
	os.Exit(run())
}

func run() int {
	exitCode := 0
	webhook := os.Getenv("FREESTUFF_WEBHOOK")
	if webhook == "" {
		log.Fatal("No webhook defined")
	}

	links, err := freestuff.GetLinksFromSubreddit("apphookup")
	if err != nil {
		log.Fatalf("Fatal error: %v", err.Error())
	}

	cache := freestuff.NewRedisCache()
	author := fmt.Sprint("Reddit Comments")

	for _, link := range links {
		known, err := cache.IsKnown(link.Link)
		if err != nil {
			exitCode = 1
			log.Printf("Cache error: %v", err)
			continue
		}

		if known || !isFree(link.Title) {
			continue
		}

		embed := discordwebhook.Embed{
			Title:     &link.Title,
			Url:       &link.Link,
			Thumbnail: &discordwebhook.Thumbnail{Url: &link.Image},
			Author: &discordwebhook.Author{
				Name: &author,
				Url:  &link.RedditLink,
			},
		}
		message := discordwebhook.Message{
			Embeds: &[]discordwebhook.Embed{embed},
		}

		time.Sleep(time.Second * 5)

		err = discordwebhook.SendMessage(webhook, message)
		if err != nil && err.Error() != "" {
			log.Printf("Discord error sending %s: %v", link.Title, err.Error())
			exitCode = 1
			continue
		}

		err = cache.SetKnown(link.Link)
		if err != nil {
			log.Printf("Cache error: %v", err.Error())
			exitCode = 1
		}
	}

	return exitCode
}

func isFree(title string) bool {
	return strings.Contains(strings.ToLower(title), "free") || strings.Contains(title, "0.00")
}
