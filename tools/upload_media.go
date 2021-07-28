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
	flag.StringVar(&host, "host", "", "host ip 或者域名")
	flag.StringVar(&secretID, "secret_id", "", "secretID")
	flag.StringVar(&secretKey, "secret_key", "", "secretKey")
	flag.IntVar(&port, "port", 80, "调用端口")
	flag.IntVar(&coroutineNum, "j", 1, "分片上传最大并行数量")
	flag.IntVar(&project, "project", 0, "TIprojectID")
	flag.IntVar(&business, "business", 0, "TIBusinessID")
	flag.StringVar(&filePath, "path", "", "要上传的文件路径")
	flag.StringVar(&mediaName, "name", "", "媒体名字")
	flag.StringVar(&mediaType, "type", "视频", "媒体类型[视频、图片、音频]")
	flag.StringVar(&mediaTag, "tag", "新闻", "媒体标签[新闻、综艺、体育、电影、电视剧、专题、互联网资讯]")
	flag.StringVar(&mediaSecondTag, "second_tag", "", "媒体二级标签[晚会、其他]")
	flag.StringVar(&mediaLang, "lang", "0", "0 普通话, 1 粤语")
	flag.Parse()
	client := media_asset_sdk.MakeMediaAssetClient(host, port, secretID, secretKey, project, business)
	mediaMeta := request.MediaMeta{
		MediaType:      mediaType,
		MediaTag:       mediaTag,
		MediaSecondTag: mediaSecondTag,
		MediaLang:      mediaLang,
	}
	if mediaName == "" {
		filenameWithSuffix := path.Base(filePath)  //获取文件名带后缀
		fileSuffix := path.Ext(filenameWithSuffix) //获取文件后缀
		mediaName = strings.TrimSuffix(filenameWithSuffix, fileSuffix)
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
	// _, err = client.DownloadToBuf(media.DownLoadURL)
	// if err != nil {
	// 	fmt.Println("DownloadToBuf failed, error: ", err)
	// 	return
	// }
	// response, reqID, err := client.DescribeCategories()
	// if err != nil {
	// 	fmt.Println("DescribeCategories failed, error: ", err, ", reqID: ", reqID)
	// 	return
	// }
	// bys, _ = json.MarshalIndent(response, "", "    ")
	// fmt.Println("DescribeCategories: ", string(bys))

	// if reqID, err := client.ModifyMedia(media.MediaID, "综艺", "晚会"); err != nil {
	// 	fmt.Println("ModifyMedia failed, error: ", err, " reqID: ", reqID)
	// 	return
	// }

	// if reqID, err := client.ModifyExpireTime(media.MediaID, 1); err != nil {
	// 	fmt.Println("ModifyExpireTime failed, error: ", err, " reqID: ", reqID)
	// 	return
	// }

	// mediaSet, tot, reqID, err := client.DescribeMedias(1, 20, &request.FilterBy{MediaNameOrID: mediaName})
	// if err != nil {
	// 	fmt.Println("DescribeMedias failed, error: ", err, " reqID: ", reqID)
	// 	return
	// }
	// bys, _ = json.MarshalIndent(mediaSet, "", "    ")
	// fmt.Println(string(bys))
	// fmt.Println("total: ", tot)
}
