package main

import (
	"fmt"
	"os"
	"time"

	imgfetch "github.com/alan-ar1/imgfetch/pkg/imgfetch"
	"github.com/charmbracelet/lipgloss"
)

func main() {

	if len(os.Args) <= 1 {
		fmt.Println("Provide an image path")
		return
	}

	imagePath := os.Args[1]

	seq, err := imgfetch.GetNativeImageSeq(imagePath)

	if err != nil {
		fmt.Println(err)
		return
	}

	info, err := getImageInfo(imagePath)
	if err != nil {
		fmt.Println(err)
		return
	}

	labelStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("4"))

	infoStr := fmt.Sprintf("%s %s\n%s %s\n%s %s\n%s %d Bytes\n%s %s\n%s %dx%d",
		labelStyle.Render("Name:"), info.name,
		labelStyle.Render("Path:"), info.absFilePath,
		labelStyle.Render("Type:"), info.fileType,
		labelStyle.Render("Size:"), info.size,
		labelStyle.Render("Modified:"), info.modTime.Format(time.DateTime),
		labelStyle.Render("Dimension:"), info.width, info.height)

	image := lipgloss.NewStyle().
		Padding(1).
		Render(seq)

	text := lipgloss.NewStyle().
		Padding(1).
		Render(infoStr)

	fmt.Println(lipgloss.JoinHorizontal(lipgloss.Top, image, text))
}
