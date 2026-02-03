package main

import (
	"flag"
	"fmt"
	"strings"
	"time"

	"github.com/alan-ar1/imgfetch/pkg/imgfetch"
	"github.com/charmbracelet/lipgloss"
)

func main() {

	infoFlag := flag.Bool("i", false, "display file info")
	flag.Parse()

	if len(flag.Args()) < 1 {
		fmt.Println("Provide an image path")
		return
	}

	filePath := flag.Args()[0]

	infoStr := ""
	imageSeq := ""
	var info FileInfo

	labelStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("4"))

	if IsHTTPURL(filePath) {
		urlFileType, err := DetectUrlFileType(filePath)
		if err != nil {
			fmt.Println(err)
			return
		}

		if strings.HasPrefix(urlFileType, "image") {
			img, _, info, err := GetImageFromURL(filePath, *infoFlag)
			if err != nil {
				fmt.Println(err)
				return
			}
			infoStr = fmt.Sprintf("%s %d Bytes\n%s %dx%d", labelStyle.Render("Size:"), info.Size, labelStyle.Render("Dimensions:"), info.Width, info.Height)
			imageSeq, err = imgfetch.GetRemoteImageSeq(img, "jpeg")
			if err != nil {
				fmt.Println(err)
				return

			}
		} else if strings.HasPrefix(urlFileType, "video") {
			thumbnailImg, err := getVideoThumbnail(filePath)
			if err != nil {
				fmt.Println(err)
				return
			}

			imageSeq, err = imgfetch.GetRemoteImageSeq(thumbnailImg, "jpeg")
			if err != nil {
				fmt.Println(err)
				return
			}

			if *infoFlag {
				info, err := GetVideoInfoFromUrl(filePath)
				if err != nil {
					fmt.Println(err)
					return
				}
				infoStr = fmt.Sprintf("%s %d Bytes\n%s %s\n",
					labelStyle.Render("Size:"), info.Size,
					labelStyle.Render("Last-Modified:"), info.LastModified)

			}

		}

	} else {
		fileType, err := DetectFileType(filePath)
		if err != nil {
			fmt.Println(err)
			return
		}
		if *infoFlag {
			info, err = GetFileInfo(filePath)
			if err != nil {
				fmt.Println(err)
				return
			}
			infoStr = fmt.Sprintf("%s %s\n%s %s\n%s %s\n%s %d Bytes\n%s %s\n",
				labelStyle.Render("Name:"), info.Name,
				labelStyle.Render("Path:"), info.AbsFilePath,
				labelStyle.Render("Type:"), fileType,
				labelStyle.Render("Size:"), info.Size,
				labelStyle.Render("Modified:"), info.ModTime.Format(time.DateTime))
		}

		if strings.HasPrefix(fileType, "image") {
			imageInfo, err := GetImageSpecInfo(filePath)
			if err != nil {
				fmt.Println(err)
			}

			infoStr += fmt.Sprintf("%s %dx%d", labelStyle.Render("Dimensions:"), imageInfo.Width, imageInfo.Height)
			imageSeq, err = imgfetch.GetImageSeq(filePath)
			if err != nil {
				fmt.Println(err)
				return
			}
		} else if strings.HasPrefix(fileType, "video") {
			thumbnailImg, err := getVideoThumbnail(filePath)
			if err != nil {
				fmt.Println(err)
				return
			}

			imageSeq, err = imgfetch.GetRemoteImageSeq(thumbnailImg, "jpeg")
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	}

	image := lipgloss.NewStyle().
		Padding(1).
		Render(imageSeq)

	text := lipgloss.NewStyle().
		Padding(1).
		Render(infoStr)

	fmt.Println(lipgloss.JoinHorizontal(lipgloss.Top, image, text))
}
