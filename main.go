package main

import "youtube_downloader/dwnld"

const dir = "PUT_YOUR_DIR_PATH_HERE"

// For example:
//const dir = "/home/lastshaman/youtube"

var urls = []string{
	"",
	"",
	"",
	"",
	"",
	"",
	"",
	"",
	"",
	"",
	"",
	"",
	"",
	"",
	"",
	"",
	"",
	"",
	"",
	"",
	"",
	"",
	"",
	"",
	"",
	"",
	"",
	"",
	"",
	"",
	"",
	"",
}

// The library used for downloading is https://github.com/kkdai/youtube
// For an example on how to download a playlist see - https://github.com/kkdai/youtube/blob/master/example_test.go.

func main() {
	dwnld.DownloadBatch(dir, urls)
}
