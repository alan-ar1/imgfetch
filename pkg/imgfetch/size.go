package imgfetch

import (
	"image"
	"os"

	"golang.org/x/sys/unix"
)

func CalculateImageTermSize(img image.Image, widthFactor int) (ImageTermSize, error) {
	w, h := img.Bounds().Dx(), img.Bounds().Dy()
	f, err := os.OpenFile("/dev/tty", unix.O_NOCTTY|unix.O_CLOEXEC|unix.O_NDELAY|unix.O_RDWR, 0666)
	if err != nil {
		return ImageTermSize{}, err
	}
	defer f.Close()

	sz, err := unix.IoctlGetWinsize(int(f.Fd()), unix.TIOCGWINSZ)
	if err != nil {
		return ImageTermSize{}, err
	}

	pixelsPerCol := float32(sz.Xpixel) / float32(sz.Col)
	pixelsPerRow := float32(sz.Ypixel) / float32(sz.Row)

	imgCols := float32(w) / pixelsPerCol
	imgRows := float32(h) / pixelsPerRow

	newImgCols := int(sz.Col) / widthFactor
	newImgRows := int(float32(newImgCols) * imgRows / imgCols)

	return ImageTermSize{Columns: newImgCols, Rows: newImgRows}, nil
}
