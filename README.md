# imgfetch
Display pixel perfect images inside terminal using kitty's graphic protocol.

<img width="2203" height="1267" alt="image" src="https://github.com/user-attachments/assets/9c427cba-492d-4131-add0-cc9d5d274c7c" />

## About
imgfetch is a Go package and a CLI tool for displaying images directly in terminal emulators that support Kitty's graphics protocol.
Also be aware that it is primarily a fun side project that I may use as a package for other TUI projects or for my personal use.

## Terminal Support
Currently only supports terminals that implement [Kitty's graphics protocol](https://sw.kovidgoyal.net/kitty/graphics-protocol/).
Other protocols like iTerm2's inline image protocol may be added in the future.
## Installation
**CLI:**
```
go install github.com/alan-ar1/imgfetch/cmd/imgfetch@latest
```
**PKG:**
```
go get github.com/alan-ar1/imgfetch/pkg/imgfetch
```

## Usage
**CLI:**
```
imgfetch path/to/image.png
```
You can also add an optional -i flag to see basic file info.\
**Note:** Video files are also supported if ffmpeg is installed. Additional file type support may be added in the future.

**PKG:**
```go
// For local files
imgSeq, err := imgfetch.GetImageSeq("local_file_path")
if err != nil {
  // handle error
}
fmt.Print(imgSeq)

// For remote files
remoteImgSeq, err := imgfetch.GetRemoteImageSeq(img) // param img of type image.Image
if err != nil {
  // handle error
}
fmt.Print(imgSeq)
```
Both `GetImageSeq` and `GetRemoteImageSeq` accept an optional second parameter to specify image dimensions:
```go
size := imgfetch.ImageTermSize{Rows: 20, Columns: 40}
imgSeq, err := imgfetch.GetImageSeq("path/to/image.png", size)
```
If not specified, images default to 1/4 of the terminal's column width.
## TODO
- [ ] iTerm2 inline image protocol support
- [ ] Image display modes (contain, cover, fill, etc.)
- [ ] PDF image preview
