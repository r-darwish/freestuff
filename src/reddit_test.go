package freestuff

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetLinksFromSubreddit(t *testing.T) {
	result, err := GetLinksFromSubreddit("apphookup")
	assert.NoError(t, err)
	assert.NotEmpty(t, result)

}
