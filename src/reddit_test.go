package freestuff

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetLinksFromSubreddit(t *testing.T) {
	result, err := GetLinksFromSubreddit("apphookup")
	assert.NoError(t, err)
	assert.NotEmpty(t, result)
}
