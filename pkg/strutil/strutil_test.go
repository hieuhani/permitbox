package strutil

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSnakeToCamel(t *testing.T) {
	res := SnakeToCamel("a_quick_brown_fox")
	assert.Equal(t, "aQuickBrownFox", res)
}

func TestSlugify(t *testing.T) {
	res := Slugify("jUmps Over_the lazy dog")
	assert.Equal(t, "jumps-over_the-lazy-dog", res)
}

func TestFormatInt(t *testing.T) {
	testSource := []any{64, "64", uint(64), int64(64)}
	for _, v := range testSource {
		val, err := FormatInt(v)
		assert.Nil(t, err)
		assert.Equal(t, "64", val)
	}
}

func TestFormatFloat(t *testing.T) {
	f := -0.7665
	assert.Equal(t, "-0.77", FormatFloat(f, 2))
}

func TestBytesToString(t *testing.T) {
	b := []byte("noice")
	assert.Equal(t, "noice", BytesToString(b))
}

func TestStringToBytes(t *testing.T) {
	s := "duc tay be"
	assert.Equal(t, []byte("duc tay be"), StringToBytes(s))
}
