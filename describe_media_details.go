package media_asset_sdk

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/Tencent-Ti/ti-sign-go/tisign"
	"github.com/Tencent-media-asset-system-sdk/media-asset-go-sdk/common"
	"github.com/Tencent-media-asset-system-sdk/media-asset-go-sdk/media_asset_service"
	"github.com/Tencent-media-asset-system-sdk/media-asset-go-sdk/media_asset_service/request"
	"github.com/Tencent-media-asset-system-sdk/media-asset-go-sdk/media_asset_service/response"
)

// DescribeMediaDetails 获取指定媒体集的详情
func (m *MediaAssetClient) DescribeMediaDetails(mediaIDs []uint64) (
	mediaSet []*response.MediaInfo, requestID string, err error) {
	if m.Port == 0 {
		m.Port = 80
	}
	action := "DescribeMediaDetails"
	service := "app-cdn4aowk"
	version := "2021-02-26"
	headerContent := tisign.HttpHeaderContent{
		XTCAction:   action,             // 请求接口
		XTCService:  service,            // 接口所属服务名
		XTCVersion:  version,            // 接口版本
		ContentType: "application/json", // http请求的content-type, 当前网关只支持: application/json  multipart/form-data
		HttpMethod:  "POST",             // http请求方法，当前网关只支持: POST GET
		Host:        m.Host,             // 访问网关的host
	}
	uri := ""
	header := map[string]string{}
	req := &request.DescribeMediaDetailsRequest{}
	req.TIBusinessID = uint32(m.TIBusinessID)
	req.TIProjectID = uint32(m.TIProjectID)
	req.MediaIDSet = mediaIDs
	if m.Inner {
		req.RequestID = common.GenerateRandomString(32)
		req.Uin = m.InnerUserName
		req.SubAccountUin = m.InnerUserName
		uri = m.InnerMediaAssetEndPoint + "/" + action
		header = nil
	} else {
		uri = fmt.Sprintf("http://%s:%d/gateway", m.Host, m.Port)
		ts := tisign.NewTiSign(headerContent, m.SecretID, m.SecretKey)
		header, _ = ts.CreateSignatureInfo()
	}
	maxTry := 3
	rsp := &response.DescribeMediaDetailsResponse{}
	for i := 0; i < maxTry; i++ {
		err = media_asset_service.HttpPost(uri, header, req, rsp)
		if rsp.Response.ApiError != nil {
			bys, _ := json.Marshal(rsp)
			err = errors.New("Response error: " + string(bys))
		}
		if err == nil {
			break
		}
	}
	return rsp.Response.MediaInfoSet, rsp.Response.RequestID, err
}
