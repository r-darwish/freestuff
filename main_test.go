package main

import (
	freestuff "github.com/r-darwish/freestuff/src"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_shouldPublish(t *testing.T) {
	freestuff.Config.PriceThreshold = 4

	type args struct {
		link freestuff.RedditLink
	}
	tests := []struct {
		name  string
		args  args
		want  freestuff.ExtraInfo
		want1 error
		want2 bool
	}{
		{"bad category", args{freestuff.RedditLink{
			Title: "[iOS] [Birthday Reminder App & Widget] [Widget Studio Pro IAP $4.99-> Free] [Keep track of upcoming birthdays, can import birthdays directly from your contacts, supports Siri Shortcuts and iCloud sync]",
			Link:  "https://apps.apple.com/app/birthday-reminder-app-widget/id1550516996"}}, nil, nil, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, got2 := shouldPublish(tt.args.link)
			assert.Equalf(t, tt.want, got, "shouldPublish(%v, %v)", tt.args.link)
			assert.Equalf(t, tt.want1, got1, "shouldPublish(%v, %v)", tt.args.link)
			assert.Equalf(t, tt.want2, got2, "shouldPublish(%v, %v)", tt.args.link)
		})
	}
}
