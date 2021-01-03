package dwnld

//import (
//	"net/http"
//	"net/url"
//
//	"github.com/kkdai/youtube/v2"
//	"google.golang.org/api/googleapi/transport"
//	youtube "google.golang.org/api/youtube/v3"
//)
//
//// downloadPlaylist downloads whole playlist (not just 1 video).
//// ytdl library doesn't support playlist-downloading out of the box, so this is a workaround suggested in https://github.com/kkdai/youtube/v2/issues/5
//func downloadPlaylist() {
//	client := &http.Client{Transport: &transport.APIKey{Key: youtubeApiKey}}
//	service, err := youtube.New(client)
//	if err != nil {
//		return // failed to connect to google api
//	}
//
//	url, err := url.Parse(playlistLink)
//	if err != nil {
//		// malformed url
//		return
//	}
//	call := service.PlaylistItems.List("id,snippet").PlaylistId(url.Query().Get("list")).MaxResults(50)
//	response, err := call.Do()
//	if err != nil {
//		return // failed to find playlist
//	}
//
//	results := response.PageInfo.TotalResults // save total results
//	err = call.Pages(nil, func(response *youtube.PlaylistItemListResponse) error {
//		for _, vid := range response.Items {
//			// query youtube video
//			md, err := ytdl.GetVideoInfo("https://youtube.com/watch?v=" + vid.Snippet.ResourceId.VideoId)
//			if err != nil {
//				// ytdl error, either
//				continue
//				// or
//				// return err
//			}
//
//			results--
//		}
//
//		return nil
//	})
//
//	// results holds the number of failed queries
//}
