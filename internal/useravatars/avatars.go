package useravatars

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
)

// Profile carries the avatar URL field for a user-facing profile.
type Profile struct {
	UserID    string `json:"user_id"`
	AvatarURL string `json:"avatar_url"`
}

// AvatarService generates avatar URLs for uploaded image data.
type AvatarService struct {
	baseURL string
}

// NewAvatarService creates an avatar service with the provided public base URL.
func NewAvatarService(baseURL string) *AvatarService {
	baseURL = strings.TrimRight(strings.TrimSpace(baseURL), "/")
	if baseURL == "" {
		baseURL = "/avatars"
	}
	return &AvatarService{baseURL: baseURL}
}

// Upload validates avatar data and returns a generated avatar_url.
func (s *AvatarService) Upload(ctx context.Context, userID string, data []byte) (string, error) {
	if err := ctx.Err(); err != nil {
		return "", err
	}
	userID = strings.TrimSpace(userID)
	if userID == "" {
		return "", errors.New("useravatars: user id is required")
	}
	if len(data) == 0 {
		return "", errors.New("useravatars: avatar data is required")
	}
	sum := sha256.Sum256(append([]byte(userID+":"), data...))
	return fmt.Sprintf("%s/%s/%s.png", s.baseURL, userID, hex.EncodeToString(sum[:8])), nil
}

// SetAvatarURL records an avatar URL on a profile.
func SetAvatarURL(profile *Profile, avatarURL string) error {
	if profile == nil {
		return errors.New("useravatars: profile is required")
	}
	avatarURL = strings.TrimSpace(avatarURL)
	if avatarURL == "" {
		return errors.New("useravatars: avatar url is required")
	}
	profile.AvatarURL = avatarURL
	return nil
}
