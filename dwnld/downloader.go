package dwnld

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/LasTshaMAN/Go-Execute/jobs"

	"github.com/rylio/ytdl"
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
	executor := jobs.NewExecutor(workersAmount, workersAmount)
	out := make(chan struct{}, workersAmount)

	for _, url := range urls {
		executor.Enqueue(func(dir, url string) func() {
			return func() {
				download(dir, url)
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

func download(dir string, url string) {
	fileName, err := downloadToTmp(dir, url)
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

func downloadToTmp(dir string, url string) (string, error) {
	fmt.Printf("downloading ... %s ... \n", url)

	vid, err := ytdl.GetVideoInfo(url)
	if err != nil {
		return "", fmt.Errorf("failed to get url info: %v", err)
	}

	assets := vid.Formats.Best(ytdl.FormatAudioBitrateKey)

	for _, asset := range assets {
		fileName, err := downloadAsset(dir, vid, asset)
		if err != nil {
			fmt.Printf("failed downloading asset: %d - error: %v \n", asset.Itag, err)
			continue
		}
		return fileName, nil
	}
	return "", fmt.Errorf("failed to download any of the assets")
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

func downloadAsset(dir string, vid *ytdl.VideoInfo, asset ytdl.Format) (string, error) {
	fileName := vid.ID + "." + asset.Extension
	if alreadyExists(dir, fileName) {
		return "", nil
	}

	file, err := os.Create(filepath.Join(tmpDirIn(dir), fileName))
	if err != nil {
		return "", fmt.Errorf("failed to create file: `%s` - error: %v", fileName, err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			fmt.Printf("failed to save video file: %s - error: %v\n", fileName, err)
		}
	}()

	if err := vid.Download(asset, file); err != nil {
		return "", fmt.Errorf("failed downloading video: %v", err)
	}

	return fileName, nil
}
