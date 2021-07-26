package main

import (
	"flag"
	"fmt"

	media_asset_sdk "github.com/Tencent-media-asset-system-sdk/media-asset-go-sdk"
	"github.com/Tencent-media-asset-system-sdk/media-asset-go-sdk/media_asset_service/request"
)

func main() {
	var host, secretID, secretKey string
	var coroutineNum, port, project, business int
	var filePath, mediaName, mediaType, mediaTag, mediaSecondTag, mediaLang string
	flag.StringVar(&host, "host", "", "host ip 或者域名")
	flag.StringVar(&secretID, "secret_id", "", "secretID")
	flag.StringVar(&secretKey, "secret_key", "", "secretKey")
	flag.IntVar(&port, "port", 80, "调用端口")
	flag.IntVar(&coroutineNum, "j", 1, "分片上传最大并行数量，默认1")
	flag.IntVar(&project, "peoject", 0, "TIpeojectID")
	flag.IntVar(&business, "business", 0, "TIBusinessID")
	flag.StringVar(&filePath, "path", "", "要上传的文件路径")
	flag.StringVar(&mediaName, "name", "", "媒体名字")
	flag.StringVar(&mediaType, "type", "", "媒体类型[视频、图片、音频]")
	flag.StringVar(&mediaTag, "tag", "新闻", "媒体标签[新闻、综艺、体育、电影、电视剧、专题、互联网资讯], 默认新闻")
	flag.StringVar(&mediaSecondTag, "second_tag", "", "媒体二级标签[晚会、其他]")
	flag.StringVar(&mediaLang, "lang", "0", "0 普通话, 1 粤语, 默认 0")
	flag.Parse()

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
	fmt.Println("Upload success, media: ", media)
}
