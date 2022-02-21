package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"path"
	"strings"

	mediaassetsdk "github.com/Tencent-media-asset-system-sdk/media-asset-go-sdk"
	"github.com/Tencent-media-asset-system-sdk/media-asset-go-sdk/media_asset_service/request"
)

func main() {
	var host, secretID, secretKey string
	var coroutineNum, port int
	var project, business uint64
	var filePath, mediaName string
	var mediaType, mediaTag, mediaSecondTag, mediaLang int
	flag.StringVar(&host, "host", "", "host ip 或者域名")
	flag.StringVar(&secretID, "secret_id", "", "secretID")
	flag.StringVar(&secretKey, "secret_key", "", "secretKey")
	flag.IntVar(&port, "port", 80, "调用端口")
	flag.IntVar(&coroutineNum, "j", 1, "分片上传最大并行数量")
	flag.Uint64Var(&project, "project", 0, "TIprojectID")
	flag.Uint64Var(&business, "business", 0, "TIBusinessID")
	flag.StringVar(&filePath, "path", "", "要上传的文件路径")
	flag.StringVar(&mediaName, "name", "", "媒体名字")
	flag.IntVar(&mediaType, "type", 2, "媒体类型")
	flag.IntVar(&mediaTag, "tag", 1, "媒体标签[新闻、综艺、体育、电影、电视剧、专题、互联网资讯]")
	flag.IntVar(&mediaSecondTag, "second_tag", 0, "媒体二级标签[晚会、其他]")
	flag.IntVar(&mediaLang, "lang", 1, "普通话 or 粤语")
	flag.Parse()
	client := mediaassetsdk.MakeMediaAssetClient(host, port, secretID, secretKey, uint32(project), uint32(business))

	// client.Inner = true
	// client.InnerUin = "1-1"
	// client.InnerInnerSubAccountUin = "superadmin"
	// client.InnerMediaAssetEndPoint = "http://media-asset-service.ai-media.svc.cluster.local:8765"
	// client.InnerDataDir = "/data/ti-platform-fs/ti-file-server"
	// client.InnerFileManagerEndPoint = "http://ti-file-manager.ti-base.svc.cluster.local:55325"
	// client.InnerFileStaticEndPoint = "http://ti-static-file-server.ti-base.svc.cluster.local:55326"

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

	// mediaSet, tot, reqID, err := client.DescribeMedias(1, 20, &request.FilterBy{MediaNameOrID: mediaName})
	// if err != nil {
	// 	fmt.Println("DescribeMedias failed, error: ", err, " reqID: ", reqID)
	// 	return
	// }
	// bys, _ = json.MarshalIndent(mediaSet, "", "    ")
	// fmt.Println(string(bys))
	// fmt.Println("total: ", tot)
	// if err := client.DownloadFile(media.DownLoadURL, "./", "temp.out"); err != nil {
	// 	fmt.Println("DownloadFile failed, error: ", err)
	// 	return
	// }
	// buf, err := client.DownloadToBuf(media.DownLoadURL)
	// if err != nil {
	// 	fmt.Println("DownloadToBuf failed, error: ", err)
	// 	return
	// }
	// fmt.Println(string(buf))
	// response, reqID, err := client.DescribeCategories()
	// if err != nil {
	// 	fmt.Println("DescribeCategories failed, error: ", err, ", reqID: ", reqID)
	// 	return
	// }
	// bys, _ = json.MarshalIndent(response, "", "    ")
	// fmt.Println("DescribeCategories: ", string(bys))

	// if reqID, err := client.ModifyMedia(media.MediaID, 2, 1); err != nil {
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

	// req := request.CreateMediasRequest{}
	// req.UploadMediaSet = append(req.UploadMediaSet, request.UploadMedia{
	// 	Name:     "质检直播流测试",
	// 	MediaURL: "https://cph-p2p-msl.akamaized.net/hls/live/2000341/test/master.m3u8",
	// 	MediaMeta: request.MediaMeta{
	// 		MediaType: 4,
	// 		MediaTag:  0,
	// 	},
	// })
	// req.UploadMediaSet = append(req.UploadMediaSet, request.UploadMedia{
	// 	Name:     "URL测试",
	// 	MediaURL: "https://ai-media-1300074211.cos.ap-shanghai.myqcloud.com/ai-media/2021-04-06/8b46057e-1923-4444-b0fb-91b094bf7530_trans.mp4",
	// 	MediaMeta: request.MediaMeta{
	// 		MediaType: 4,
	// 		MediaTag:  0,
	// 	},
	// })
	// req.UploadMediaSet = append(req.UploadMediaSet, request.UploadMedia{
	// 	Name:      "测试2",
	// 	LocalPath: "/data/test.mp4",
	// 	MediaMeta: request.MediaMeta{
	// 		MediaType: 4,
	// 		MediaTag:  0,
	// 	},
	// })
	// rsp, err := client.CreateMedias(&req)
	// if err == nil && rsp.Response.ApiError == nil {
	// 	bys, _ := json.Marshal(rsp)
	// 	fmt.Println(string(bys))
	// } else {
	// 	fmt.Print()
	// }

}
