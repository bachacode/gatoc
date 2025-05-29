package assets

import (
	"embed"
)

//go:embed images
var imagesFS embed.FS

// Optionally provide a function to get a sub-filesystem
func GetCatGif() ([]byte, error) {
	return imagesFS.ReadFile("images/cat.gif")
}
