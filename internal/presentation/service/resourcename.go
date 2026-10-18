package service

import (
	"errors"

	"go.einride.tech/aip/resourcename"
)

var (
	ErrInvalidImageResourceName = errors.New("invalid image resource name")
)

const (
	imageResourcePattern = "images/{image}"
)

func parseImageName(name string) (string, error) {
	var imageID string
	if err := resourcename.Sscan(name, imageResourcePattern, &imageID); err != nil {
		return "", ErrInvalidImageResourceName
	}

	return imageID, nil
}

func formatImageName(imageID string) string {
	return resourcename.Sprint(imageResourcePattern, imageID)
}
