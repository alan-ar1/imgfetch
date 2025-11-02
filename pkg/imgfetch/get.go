package imgfetch

import (
	"errors"
	"os"

	"github.com/alan-ar1/imgfetch/pkg/imgfetch/common"
	"github.com/alan-ar1/imgfetch/pkg/imgfetch/kitty"
)

type ImageTermSize = common.ImageTermSize

func GetNativeImageSeq(imagePath string, size ...ImageTermSize) (string, error) {

	if os.Getenv("GHOSTTY_RESOURCES_DIR") != "" {
		if size == nil {
			imageTermSize, err := GetImageTermSize(imagePath, 4)
			if err != nil {
				return "", err
			}
			size = []ImageTermSize{imageTermSize}
		}
		seq, err := kitty.Kitty(imagePath, size[0])
		if err != nil {
			return "", err
		}
		return seq, nil
	} else {
		return "", errors.New("Terminal doesn't support kitty's graphic protocol")
	}
}
