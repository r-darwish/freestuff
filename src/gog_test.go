package freestuff

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetGogInfo(t *testing.T) {
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
			args:    args{"https://www.gog.com/game/shadow_tactics_blades_of_the_shogun"},
			want:    &GogInfo{4.8},
			wantErr: assert.NoError,
		}, {
			args:    args{"https://www.gog.com/game/crime_cities"},
			want:    &GogInfo{3.9},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetGogInfo(tt.args.link)
			if !tt.wantErr(t, err, fmt.Sprintf("GetGogInfo(%v)", tt.args.link)) {
				return
			}
			assert.Equalf(t, tt.want, got, "GetGogInfo(%v)", tt.args.link)
		})
	}
}
