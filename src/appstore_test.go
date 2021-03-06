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
			want: AppstoreInfo{
				Score:    4.7,
				Ratings:  58,
				Category: "Games",
			},
			wantErr: assert.NoError,
		},
		{
			args: args{"https://apps.apple.com/app/id482361332"},
			want: AppstoreInfo{
				Score:    4.8,
				Ratings:  6900,
				Category: "Weather",
			},
			wantErr: assert.NoError,
		},
		{
			args: args{"https://apps.apple.com/app/convusic/id1591366129"},
			want: AppstoreInfo{
				Score:    0,
				Ratings:  0,
				Category: "Music",
			},
			wantErr: assert.NoError,
		},
		{
			args: args{"https://apps.apple.com/app/twoslideover/id1547137384"},
			want: AppstoreInfo{
				Score:    5,
				Ratings:  1,
				Category: "Photo & Video",
			},
			wantErr: assert.NoError,
		},
		{
			args: args{"https://apps.apple.com/us/app/easy-spending-budget/id437238261"},
			want: AppstoreInfo{
				Score:    5,
				Ratings:  1,
				Category: "Photo & Video",
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetExtraInfo(tt.args.link)
			if !tt.wantErr(t, err, fmt.Sprintf("GetAppstoreInfo(%v)", tt.args.link)) {
				return
			}
			assert.Equalf(t, tt.want, got, "GetAppstoreInfo(%v)", tt.args.link)
		})
	}
}
