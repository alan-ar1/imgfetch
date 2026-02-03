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
	label := labelStyle.Render

	if IsHTTPURL(filePath) {
		urlFileType, err := DetectUrlFileType(filePath)
		if err != nil {
			fmt.Println(err)
			return
		}

		if strings.HasPrefix(urlFileType, "image") {
			img, contentType, info, err := GetImageFromURL(filePath, *infoFlag)
			if err != nil {
				fmt.Println(err)
				return
			}

			if *infoFlag {
				infoStr = fmt.Sprintf("%s %s\n%s %d Bytes\n%s %dx%d", label("Type:"), contentType, label("Size:"), info.Size, label("Dimensions:"), info.Width, info.Height)
			}

			imageSeq, err = imgfetch.GetRemoteImageSeq(img)
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

			imageSeq, err = imgfetch.GetRemoteImageSeq(thumbnailImg)
			if err != nil {
				fmt.Println(err)
				return
			}

			if *infoFlag {
				info, contentType, err := GetVideoInfoFromUrl(filePath)
				if err != nil {
					fmt.Println(err)
					return
				}
				infoStr = fmt.Sprintf("%s %s\n%s %d Bytes\n%s %s\n",
					label("Type:"), contentType,
					label("Size:"), info.Size,
					label("Last-Modified:"), info.LastModified)

			}

		} else {
			fmt.Println("Url not supported")
			return
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
				label("Name:"), info.Name,
				label("Path:"), info.AbsFilePath,
				label("Type:"), fileType,
				label("Size:"), info.Size,
				label("Modified:"), info.ModTime.Format(time.DateTime))
		}

		if strings.HasPrefix(fileType, "image") {
			if *infoFlag {
				imageInfo, err := GetImageSpecInfo(filePath)
				if err != nil {
					fmt.Println(err)
				}
				infoStr += fmt.Sprintf("%s %dx%d", label("Dimensions:"), imageInfo.Width, imageInfo.Height)
			}
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

			imageSeq, err = imgfetch.GetRemoteImageSeq(thumbnailImg)
			if err != nil {
				fmt.Println(err)
				return
			}
		} else {
			fmt.Println("File not supported")
			return
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
