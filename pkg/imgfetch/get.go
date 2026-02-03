package imgfetch

import (
	"image"
	"os"

	"github.com/alan-ar1/imgfetch/pkg/imgfetch/common"
	"github.com/alan-ar1/imgfetch/pkg/imgfetch/kitty"
	_ "image/jpeg"
	_ "image/png"
)

type ImageTermSize = common.ImageTermSize

/*
var isKittyProtocolSupported bool = os.Getenv("GHOSTTY_RESOURCES_DIR") != "" ||
	os.Getenv("KITTY_WINDOW_ID") != "" ||
	os.Getenv("KONSOLE_VERSION") != "" ||
	os.Getenv("WARP_IS_LOCAL_SHELL") != "" ||
	os.Getenv("WAYST_VERSION") != "" ||
	os.Getenv("WEZTERM_EXECUTABLE") != ""
*/

func GetImageSeq(imagePath string, size ...ImageTermSize) (string, error) {

	if size == nil {
		file, err := os.Open(imagePath)
		if err != nil {
			return "", err
		}
		defer file.Close()
		img, _, err := image.Decode(file)

		imageTermSize, err := CalculateImageTermSize(img, 4)
		if err != nil {
			return "", err
		}
		size = []ImageTermSize{imageTermSize}
	}
	seq, err := kitty.GetSeq(imagePath, size[0])
	if err != nil {
		return "", err
	}
	return seq, nil

}

func GetRemoteImageSeq(img image.Image, format string, size ...ImageTermSize) (string, error) {

	if size == nil {
		imageTermSize, err := CalculateImageTermSize(img, 4)
		if err != nil {
			return "", err
		}
		size = []ImageTermSize{imageTermSize}
	}
	seq, err := kitty.GetRemoteSeq(img, format, size[0])
	if err != nil {
		return "", err
	}
	return seq, nil

}
