package response

import "github.com/Tencent-media-asset-system-sdk/media-asset-go-sdk/media_asset_service/request"

// DescribeCategoriesResponse
type DescribeCategoriesResponse struct {
	Response struct {
		CommonResponse
		CategorySet []*request.Category `json:"CategorySet"`
		LabelSet    []*request.Label    `json:"LabelSet"`
	} `json:"Response"`
}
