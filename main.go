package main

import (
	"fmt"
	"strings"
	"youtube_downloader/dwnld"
)

const dir = "/Users/iurii/youtube"

var urls = ""

// The library used for downloading is https://github.com/kkdai/youtube
// For an example on how to download a playlist see - https://github.com/kkdai/youtube/blob/master/example_test.go.

func main() {
	if err := dwnld.CleanTmpArtifacts(dir); err != nil {
		fmt.Printf("failed cleaning tmp artifacts: %v \n", err)
		return
	}

	dwnld.DownloadBatch(dir, strings.Split(urls, "\n"))
	//dwnld.DownloadPlaylist(dir, "PLL1JDiTNCUAVq9YeGbxDtqBgaqUZajGIH")
}
