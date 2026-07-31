package useravatars

import (
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUploadReturnsGeneratedAvatarURL(t *testing.T) {
	service := NewAvatarService("https://cdn.example.com/avatars")

	url, err := service.Upload(context.Background(), "user-1", []byte("image-bytes"))

	require.NoError(t, err)
	assert.True(t, strings.HasPrefix(url, "https://cdn.example.com/avatars/user-1/"))
	assert.True(t, strings.HasSuffix(url, ".png"))
}

func TestSetAvatarURL(t *testing.T) {
	profile := &Profile{UserID: "user-1"}

	err := SetAvatarURL(profile, "https://cdn.example.com/a.png")

	require.NoError(t, err)
	assert.Equal(t, "https://cdn.example.com/a.png", profile.AvatarURL)
}

func TestUploadRejectsEmptyData(t *testing.T) {
	service := NewAvatarService("")

	_, err := service.Upload(context.Background(), "user-1", nil)

	require.Error(t, err)
}
