package dwnld

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
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

	fmt.Printf("will download vids in `%s` \n", dir)

	workersAmount := uint(100)
	e := executor.New(workersAmount)
	out := make(chan struct{}, workersAmount)

	client := youtube.Client{
		HTTPClient: http.DefaultClient,
	}

	for _, url := range urls {
		e.Enqueue(func(dir, url string) func() {
			return func() {
				fmt.Printf("downloading ... %s ... \n", url)

				vid, err := client.GetVideoContext(context.Background(), url)
				if err != nil {
					panic(fmt.Errorf("failed to get url %s info: %v", url, err))
				}

				downloadToTmp(dir, vid)
				out <- struct{}{}
			}
		}(dir, url))
	}

	for range urls {
		<-out
	}

	fmt.Printf("done \n")
}

func DownloadPlaylist(dir string, playlistID string) {
	fmt.Printf("will download playlist in `%s` \n", dir)

	workersAmount := uint(100)
	e := executor.New(workersAmount)
	out := make(chan struct{}, workersAmount)

	client := youtube.Client{
		HTTPClient: http.DefaultClient,
	}

	playlist, err := client.GetPlaylist(playlistID)
	if err != nil {
		panic(err)
	}

	for _, vid := range playlist.Videos {
		e.Enqueue(func(dir string, vid *youtube.PlaylistEntry) func() {
			return func() {
				fmt.Printf("downloading ... %s ... \n", vid.ID)

				video, err := client.VideoFromPlaylistEntry(vid)
				if err != nil {
					panic(err)
				}

				downloadToTmp(dir, video)
				out <- struct{}{}
			}
		}(dir, vid))
	}

	for range playlist.Videos {
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

func CleanTmpArtifacts(dir string) error {
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

func downloadToTmp(dir string, vid *youtube.Video) {
	fileName, downloaded, err := tryDownload(dir, vid)
	if err != nil {
		fmt.Printf("download failed for vid: %s - error: %v \n", vid.ID, err)
		return
	}
	if !downloaded {
		return
	}

	err = os.Rename(filepath.Join(tmpDirIn(dir), fileName), filepath.Join(dir, fileName))
	if err != nil {
		fmt.Printf("failed to move file: %s - error: %v \n", fileName, err)
		return
	}
}

func tryDownload(dir string, vid *youtube.Video) (string, bool, error) {
	fileName := vid.ID + ".mp4"

	downloaded, err := maybeDownloadVid(dir, fileName, vid)
	if err != nil {
		fmt.Printf("failed downloading vid: %s - error: %v \n", vid.Title, err)
	}
	if !downloaded {
		return "", false, nil
	}

	return fileName, true, nil
}

func alreadyExists(dir, fileName string) bool {
	stats, err := os.Stat(filepath.Join(dir, fileName))
	if err != nil {
		if os.IsNotExist(err) {
			return false
		} else {
			panic(fmt.Sprintf("Failed to check whether file already exists: %s - error: %v", fileName, err))
		}
	}
	return stats.Size() != 0
}

func maybeDownloadVid(dir string, fileName string, vid *youtube.Video) (bool, error) {
	if alreadyExists(dir, fileName) {
		return false, nil
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

	// Interested only in formats with audio.
	formats := vid.Formats.WithAudioChannels()

	// Pick the first format that allows for downloading full stream (I've encountered an issue that some formats
	// stick EOF in the middle of its stream or something, thus not allowing to download the full stream, only
	// some part of it).
	for _, format := range formats {
		stream, size, err := dwd.GetStream(vid, &format)
		if err != nil {
			return false, fmt.Errorf("can't get context for vid: %s", vid.ID)
		}

		// For some reason there are formats of 0 size ... we don't need these.
		if size == 0 {
			continue
		}

		file, err := os.Create(dir + "/tmp/" + fileName)
		if err != nil {
			return false, fmt.Errorf("can't create file: %s, err: %w", file.Name(), err)
		}

		n, err := io.Copy(file, stream)
		if err != nil {
			return false, fmt.Errorf("can't copy stream for file: %s, err: %w", file.Name(), err)
		}

		err = file.Close()
		if err != nil {
			return false, fmt.Errorf("can't close file: %s, err: %w", file.Name(), err)
		}

		// Once fully downloadable stream is found - it's good for us.
		if n == size {
			break
		}
	}

	return true, nil
}
