package freestuff

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPrice(t *testing.T) {
	type args struct {
		title string
	}
	tests := []struct {
		name    string
		args    args
		want    float64
		wantErr assert.ErrorAssertionFunc
	}{
		{
			args:    args{"[Win/Mac][VideoProc Converter][$78.90>Free]"},
			want:    78.90,
			wantErr: assert.NoError,
		},
		{
			args:    args{"[iOS/iPadOS] [Guess Word! Forehead Charade] [Premium Forever Promo $9.99 -> $0.00]"},
			want:    9.99,
			wantErr: assert.NoError,
		},
		{
			args:    args{"[Android] [Batrix] [$0.99 → Free] [A Matrix live wallpaper that reacts to the battery state. Posted by Dev. No ads or IAPs. Rating 4.8]"},
			want:    0.99,
			wantErr: assert.NoError,
		},
		{
			args:    args{"[iOS][Calory: Simple Calories counter][lifetime IAP 29,99$ —-Free]"},
			want:    29.99,
			wantErr: assert.NoError,
		}, {
			args:    args{"[Bad North: Jotunn Edition [€4,49 -> €2.29]"},
			want:    4.49,
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Price(tt.args.title)
			if !tt.wantErr(t, err, fmt.Sprintf("Price(%v)", tt.args.title)) {
				return
			}
			assert.Equalf(t, tt.want, got, "Price(%v)", tt.args.title)
		})
	}
}
