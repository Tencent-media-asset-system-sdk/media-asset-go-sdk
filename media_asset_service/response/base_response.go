package response

import "github.com/Tencent-media-asset-system-sdk/media-asset-go-sdk/media_asset_service/request"

// Response
type Response interface {
	SetRequestId(requestId string)
}

// BaseResponse
type BaseResponse struct {
	Response interface{} `json:"Response"`
}

// CommonResponse
type CommonResponse struct {
	RequestID string `json:"RequestID,omitempty"`
	ApiError  *Error `json:"Error,omitempty"`
	Error     error  `json:"-"`
}

// Error
type Error struct {
	Code    string `json:"Code"`
	Message string `json:"Message"`
}

// MediaInfo
type MediaInfo struct {
	MediaID  uint64 `json:"MediaID"`
	Name     string `json:"Name"`
	Duration uint32 `json:"Duration"`
	Size     uint32 `json:"Size"`
	Width    uint32 `json:"Width"`
	Height   uint32 `json:"Height"`
	FPS      uint32 `json:"FPS"`
	BitRate  uint32 `json:"BitRate"`
	Format   string `json:"Format"`
	request.MediaMeta
	DownLoadURL  string `json:"DownLoadURL"`
	Status       string `json:"Status"`
	FailedReason string `json:"FailedReason"`
	Bucket       string  `json:"Bucket,omitempty"`
	Key          string  `json:"Key,,omitempty"`
	UploadId     string  `json:"UploadId,omitempty"`
	LocalPath    string  `json:"LocalPath,omitempty"`
}

// FailedMediaInfo
type FailedMediaInfo struct {
	MediaID      uint64 `json:"MediaID"`
	FailedReason string `json:"FailedReason"`
}
