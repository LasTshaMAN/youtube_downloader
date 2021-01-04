package dwnld

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/LasTshaMAN/Go-Execute/executor"
	"github.com/kkdai/youtube/v2"
	"github.com/kkdai/youtube/v2/downloader"
)

const tmpDirName = "tmp"

func DownloadBatch(dir string, urls []string) {
	urls = removeEmpty(urls)
	urls = removeDuplicates(urls)

	if err := cleanTmpArtifacts(dir); err != nil {
		fmt.Printf("failed cleaning tmp artifacts: %v \n", err)
		return
	}

	fmt.Printf("will download in `%s` \n", dir)

	workersAmount := uint(100)
	e := executor.New(workersAmount)
	out := make(chan struct{}, workersAmount)

	client := youtube.Client{
		HTTPClient: http.DefaultClient,
	}

	for _, url := range urls {
		e.Enqueue(func(dir, url string) func() {
			return func() {
				download(client, dir, url)
				out <- struct{}{}
			}
		}(dir, url))
	}

	for range urls {
		<-out
	}

	fmt.Printf("done \n")
}

func removeEmpty(src []string) (result []string) {
	for _, str := range src {
		if str != "" {
			result = append(result, str)
		}
	}
	return
}

func removeDuplicates(src []string) []string {
	m := make(map[string]struct{})
	for _, item := range src {
		m[item] = struct{}{}
	}
	result := make([]string, 0, len(m))
	for item := range m {
		result = append(result, item)
	}
	return result
}

func cleanTmpArtifacts(dir string) error {
	tmpDirPath := tmpDirIn(dir)
	if err := os.RemoveAll(tmpDirPath); err != nil {
		return err
	}
	if err := os.MkdirAll(tmpDirPath, os.ModePerm); err != nil {
		panic(fmt.Sprintf("Failed to create tmp dir: %v", err))
	}
	return nil
}

func tmpDirIn(dir string) string {
	return filepath.Join(dir, tmpDirName)
}

func download(client youtube.Client, dir string, url string) {
	fileName, err := downloadToTmp(client, dir, url)
	if err != nil {
		fmt.Printf("download failed: %s - error: %v \n", url, err)
		return
	}
	if fileName == "" {
		return
	}
	if err := os.Rename(filepath.Join(tmpDirIn(dir), fileName), filepath.Join(dir, fileName)); err != nil {
		fmt.Printf("failed to move file: %s - error: %v \n", fileName, err)
		return
	}
}

func downloadToTmp(client youtube.Client, dir string, url string) (string, error) {
	fmt.Printf("downloading ... %s ... \n", url)

	vid, err := client.GetVideoContext(context.Background(), url)
	if err != nil {
		return "", fmt.Errorf("failed to get url info: %v", err)
	}

	//format := vid.Formats.FindByQuality("hd1080")

	fileName := vid.ID + ".mp4"

	err = downloadVid(dir, fileName, vid)
	if err != nil {
		fmt.Printf("failed downloading vid: %s - error: %v \n", vid.Title, err)
	}

	return fileName, nil
}

func alreadyExists(dir, fileName string) bool {
	if _, err := os.Stat(filepath.Join(dir, fileName)); err != nil {
		if os.IsNotExist(err) {
			return false
		} else {
			panic(fmt.Sprintf("Failed to check whether file already exists: %s - error: %v", fileName, err))
		}
	}
	return true
}

func downloadVid(dir string, fileName string, vid *youtube.Video) error {
	if alreadyExists(dir, fileName) {
		return nil
	}

	httpTransport := &http.Transport{
		Proxy:                 http.ProxyFromEnvironment,
		IdleConnTimeout:       60 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

	dwd := &downloader.Downloader{
		OutputDir: dir + "/tmp",
	}
	dwd.HTTPClient = &http.Client{Transport: httpTransport}

	err := dwd.DownloadWithHighQuality(context.Background(), fileName, vid, "hd1080")
	if err == nil {
		return nil
	}

	return dwd.Download(context.Background(), vid, &vid.Formats[0], fileName)
}
