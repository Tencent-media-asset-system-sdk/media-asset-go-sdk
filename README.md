# media-asset-go-sdk

此SDK用于在go语言中中方便的向媒体管理系统上传媒体资源。使用此SDK之前请先参考相关的API接口文档。


## 构建客户端

所有的方法都封装在 `mediaassetsdk.MediaAssetClient`。构造一个`MediaAssetClient`的方法如下

```go
import (
  mediaassetsdk "github.com/Tencent-media-asset-system-sdk/media-asset-go-sdk"
)

host := "106.52.71.124" // 调用服务的的 host
port := "" // 调用服务的的 port, 80 填空
secretID := "secretID" // secretID
secretID := "secretKey" // secretKey
businessID := 1 // 业务ID
projectID := 1 // 项目ID
// TODO: config below
client := mediaassetsdk.MakeMediaAssetClient(host, port, secretID, secretKey, projectID, businessID)
```

## 获取支持媒体列表
```go
response, reqID, err := client.DescribeCategories()
// response 返回的支持的媒体列表信息 *response.DescribeCategoriesResponse 结构
if err != nil {
  fmt.Println("DescribeCategories failed, error: ", err, ", reqID: ", reqID)
} else {
  bys, _ = json.MarshalIndent(response, "", "    ")
  fmt.Println("DescribeCategories: ", string(bys))
}
```

## 上传媒体
```go
import "github.com/Tencent-media-asset-system-sdk/media-asset-go-sdk/media_asset_service/request"

filePath := "data/test.mp4" // 文件路径
mediaName := "test视频" // 媒体名称
mediaType := "视频" // 媒体类型, 可选 视频，音频，图片
mediaTag := "新闻" // 媒体标签, 可选 新闻, 综艺, 电影, 电视剧, 体育, 专题, 互联网资讯
mediaSecondTag := "" // 二级标签，如果一级标签为综艺可选 "晚会" 和 "其他"，其他为空
mediaLang := "普通话" // 可选 普通话, 粤语
coroutineNum := 4 // 并发上传的协程数
mediaMeta := request.MediaMeta{
  MediaType:      mediaType,
  MediaTag:       mediaTag,
  MediaSecondTag: mediaSecondTag,
  MediaLang:      mediaLang,
}

// 上传文件到媒体管理系统
media, reqSet, err := client.UploadFile(filePath, mediaName, mediaMeta, coroutineNum)
// media 上传成功后的媒体信息 *response.MediaInfo 结构
if err != nil {
  fmt.Println("Upload failed, error: ", err, " RequestIDSet: ", reqSet)
} else {
  bys, _ := json.MarshalIndent(media, "", "    ")
  fmt.Println("Upload success, media: ", string(bys))
}

// 上传内存到媒体管理系统
filebuf, err := ioutil.ReadFile(filePath)
if err != nil {
  media, reqSet, err := client.UploadBuf(filebuf, mediaName, mediaMeta, coroutineNum)
} else {
  fmt.Println("Read file error: ", err)
}
```

## 获取指定媒体详细信息
```go
medias, reqID, err := client.DescribeMediaDetails([]uint64{media.MediaID})
// medias 返回的媒体信息列表 []*response.MediaInfo
```

## 下载媒体
```go
// 下载媒体到文件
dir := "./data" // 下载到的目录
fileName := "download.mp4" // 下载的文件名
if err := client.DownloadFile(media.DownLoadURL, dir, fileName); err != nil {
  fmt.Println("DownloadFile failed, error: ", err)
}

// 下载媒体到内存
buf, err := client.DownloadToBuf(media.DownLoadURL)
if err != nil {
  fmt.Println("DownloadToBuf failed, error: ", err)
}
```

## 获取上传媒体列表
```go
medias, tot, reqID, err := client.DescribeMedias(1, 20, &request.FilterBy{MediaNameOrID: mediaName})
// medias 返回的媒体列表 []*response.MediaInfo
// tot 列表的总媒体数
if err != nil {
  fmt.Println("DescribeMedias failed, error: ", err, " reqID: ", reqID)
}
```

## 删除媒体
```go
failedMedias, reqID, err := client.RemoveMedias([]uint64{media.MediaID})
// failedMedias 删除失败的媒体信息列表 []*response.FailedMediaInfo 结构
```

## 修改媒体类型
```go
newTag := ""综艺" // 新标签
newSeconeTag := "晚会" // 新二级标签
reqID, err := client.ModifyMedia(media.MediaID, newTag, newSeconeTag);
```

## 修改媒体过期时间
```go
day := 1 // 媒体过期时间
reqID, err := client.ModifyExpireTime(media.MediaID, day)
```

## 批量创建媒体
```go
req := request.CreateMediasRequest{}
req.UploadMediaSet = append(req.UploadMediaSet, request.UploadMedia{
  Name:     "直播流测试1",
  MediaURL: "http://live.tencent.com",
  MediaMeta: request.MediaMeta{
    MediaType: "直播流",
    MediaTag:  "新闻",
  },
})
req.UploadMediaSet = append(req.UploadMediaSet, request.UploadMedia{
  Name:     "测试2",
  MediaURL: "http://video.tencent.com",
  MediaMeta: request.MediaMeta{
    MediaType: "视频",
    MediaTag:  "新闻",
  },
})
req.UploadMediaSet = append(req.UploadMediaSet, request.UploadMedia{
  Name:     "测试2",
  LocalPath: "/data/test.mp4",
  MediaMeta: request.MediaMeta{
    MediaType: "视频",
    MediaTag:  "新闻",
  },
})
rsp, err := client.CreateMedias(&req)
if err == nil && rsp.Response.ApiError == nil {
  bys, _ := json.Marshal(rsp)
  fmt.Println(string(bys))
}
```

## 工具
上传单个媒体工具

go run
```
go run tools/upload_media.go -h
```

centos
```
./tools/upload_media -h
```
