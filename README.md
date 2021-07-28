# media-asset-go-sdk
媒体AI中台媒体管理系统SDK

# API LIST
```go
// UploadFile 上传文件
// filePath 文件路径
// mediaName 媒体名称
// coroutineNum 上传最大并发协程数
// mediaInfo request.MediaMeta 媒体的类型和标签信息
UploadFile(filePath, mediaName string, mediaMeta request.MediaMeta, coroutineNum int) (
  media *response.MediaInfo, requestIDSet []string, err error)

// DownloadFile 通过媒体信息返回的url下载文件到本地
DownloadFile(downloadURL, dir, fileName string) (err error)

// DownloadToBuf 通过媒体信息返回的url下载文件到内存
DownloadToBuf(downloadURL string) (buf []byte, err error)

// DescribeMedias 拉取媒体列表
DescribeMedias(pageNumber, pageSize int, filterBy *request.FilterBy) (
  mediaSet []*response.MediaInfo, TotalCount int, requestID string, err error)

// DescribeMediaDetails 获取指定媒体集的详情
DescribeMediaDetails(mediaIDs []uint64) (mediaSet []*response.MediaInfo, requestID string, err error)

// RemoveMedias 删除指定媒体集
RemoveMedias(mediaIDs []uint64) (failedMediaSet []*response.FailedMediaInfo, requestID string, err error)

// DescribeCategories 返回可选媒体类型列表
DescribeCategories() (categortSet *response.DescribeCategoriesResponse, requestID string, err error)

// ModifyMedia 修改媒体信息
ModifyMedia(mediaID uint64, mediaTag, mediaSecondTag string) (requestID string, err error)

// ModifyExpireTime 修改文件过期时间，当前时间算起来，有效时间为 days 天
ModifyExpireTime(mediaID uint64, days int) (requestID string, err error)
```

# USAGE
```go
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"path"
	"strings"

	media_asset_sdk "github.com/Tencent-media-asset-system-sdk/media-asset-go-sdk"
	"github.com/Tencent-media-asset-system-sdk/media-asset-go-sdk/media_asset_service/request"
)

func main() {
	var host, secretID, secretKey string
	var coroutineNum, port, project, business int
	var filePath, mediaName, mediaType, mediaTag, mediaSecondTag, mediaLang string
  // TODO: config
	client := media_asset_sdk.MakeMediaAssetClient(host, port, secretID, secretKey, project, business)
	mediaMeta := request.MediaMeta{
		MediaType:      mediaType,
		MediaTag:       mediaTag,
		MediaSecondTag: mediaSecondTag,
		MediaLang:      mediaLang,
	}
	media, reqSet, err := client.UploadFile(filePath, mediaName, mediaMeta, coroutineNum)
	if err != nil {
		fmt.Println("Upload failed, error: ", err, " RequestIDSet: ", reqSet)
		return
	}
	bys, _ := json.MarshalIndent(media, "", "    ")
	fmt.Println("Upload success, media: ", string(bys))

	if err := client.DownloadFile(media.DownLoadURL, "./", "temp.out"); err != nil {
		fmt.Println("DownloadFile failed, error: ", err)
		return
	}
	_, err = client.DownloadToBuf(media.DownLoadURL)
	if err != nil {
		fmt.Println("DownloadToBuf failed, error: ", err)
		return
	}
	response, reqID, err := client.DescribeCategories()
	if err != nil {
		fmt.Println("DescribeCategories failed, error: ", err, ", reqID: ", reqID)
		return
	}
	bys, _ = json.MarshalIndent(response, "", "    ")
	fmt.Println("DescribeCategories: ", string(bys))

	if reqID, err := client.ModifyMedia(media.MediaID, "综艺", "晚会"); err != nil {
		fmt.Println("ModifyMedia failed, error: ", err, " reqID: ", reqID)
		return
	}

	if reqID, err := client.ModifyExpireTime(media.MediaID, 1); err != nil {
		fmt.Println("ModifyExpireTime failed, error: ", err, " reqID: ", reqID)
		return
	}

	mediaSet, tot, reqID, err := client.DescribeMedias(1, 20, &request.FilterBy{MediaNameOrID: mediaName})
	if err != nil {
		fmt.Println("DescribeMedias failed, error: ", err, " reqID: ", reqID)
		return
	}
	bys, _ = json.MarshalIndent(mediaSet, "", "    ")
	fmt.Println(string(bys))
	fmt.Println("total: ", tot)
}
```

# Tool
上传单个视频工具

go run
```
go run tools/upload_media.go -h
```

centos
```
./tools/upload_media -h
```