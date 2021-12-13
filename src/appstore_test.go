package freestuff

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetAppstoreInfo(t *testing.T) {
	type args struct {
		link string
	}
	tests := []struct {
		name    string
		args    args
		want    ExtraInfo
		wantErr assert.ErrorAssertionFunc
	}{
		{
			args: args{"https://apps.apple.com/app/demetrios/id1295749833"},
			want: &AppstoreInfo{
				Score:   4.8,
				Ratings: 55,
			},
			wantErr: assert.NoError,
		},
		{
			args: args{"https://apps.apple.com/app/id482361332"},
			want: &AppstoreInfo{
				Score:   4.8,
				Ratings: 6800,
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetAppstoreInfo(tt.args.link)
			if !tt.wantErr(t, err, fmt.Sprintf("GetAppstoreInfo(%v)", tt.args.link)) {
				return
			}
			assert.Equalf(t, tt.want, got, "GetAppstoreInfo(%v)", tt.args.link)
		})
	}
}
