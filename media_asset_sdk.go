package mediaassetsdk

import (
	"github.com/Tencent-media-asset-system-sdk/media-asset-go-sdk/media_asset_service/request"
	"github.com/Tencent-media-asset-system-sdk/media-asset-go-sdk/media_asset_service/response"
)

const (
	SERVICE = "app-cdn4aowk"
	VERSION = "2021-02-26"
)

// MediaAssetFunction 媒体管理系统sdk功能列表
type MediaAssetFunction interface {
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
		mediaSet []*response.MediaInfo, totalCount int, requestID string, err error)

	// DescribeMediaDetails 获取指定媒体集的详情
	DescribeMediaDetails(mediaIDs []string) (mediaSet []*response.MediaInfo, requestID string, err error)

	// RemoveMedias 删除指定媒体集
	RemoveMedias(mediaIDs []string) (failedMediaSet []*response.FailedMediaInfo, requestID string, err error)

	// DescribeCategories 返回可选媒体类型列表
	DescribeCategories() (categortSet *response.DescribeCategoriesResponse, requestID string, err error)

	// ModifyMedia 修改媒体信息
	ModifyMedia(mediaID string, mediaTag, mediaSecondTag string) (requestID string, err error)

	// ModifyExpireTime 修改文件过期时间，当前时间算起来，有效时间为 days 天
	ModifyExpireTime(mediaID string, days int) (requestID string, err error)

	// CreateMedias 批量创建媒体
	CreateMedias(req *request.CreateMediasRequest) (rsp *response.CreateMediasResponse, err error)
}

// sdk客户端
type MediaAssetClient struct {
	Host                     string
	Port                     int
	SecretID                 string
	SecretKey                string
	TIProjectID              uint32
	TIBusinessID             uint32
	Inner                    bool
	InnerMediaAssetEndPoint  string
	InnerFileManagerEndPoint string
	InnerFileStaticEndPoint  string
	InnerUserName            string
	InnerDataDir             string
}

// MakeMediaAssetClient 创建一个客户端
func MakeMediaAssetClient(host string, port int, secretID, secretKey string,
	tiProjectID, tiBusinessID uint32) *MediaAssetClient {
	if port == 0 {
		port = 80
	}
	return &MediaAssetClient{
		Host:         host,
		Port:         port,
		SecretID:     secretID,
		SecretKey:    secretKey,
		TIProjectID:  tiProjectID,
		TIBusinessID: tiBusinessID,
		Inner:        false,
	}
}

// CheckStatusFailed 检查媒体状态是否是上传失败
func (client MediaAssetClient) CheckStatusFailed(status int) bool {
	return status == 0 || status == 4 || status == 7 || status == 9 || status == 12 || status == 13
}

// CheckStatusSuccess 检查媒体状态是否是上传成功
func (client MediaAssetClient) CheckStatusSuccess(status int) bool {
	return status == 8
}
