package mediaassetsdk

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/Tencent-Ti/ti-sign-go/tisign"
	"github.com/Tencent-media-asset-system-sdk/media-asset-go-sdk/common"
	mediaassetservice "github.com/Tencent-media-asset-system-sdk/media-asset-go-sdk/media_asset_service"
	"github.com/Tencent-media-asset-system-sdk/media-asset-go-sdk/media_asset_service/request"
	"github.com/Tencent-media-asset-system-sdk/media-asset-go-sdk/media_asset_service/response"
)

// ModifyExpireTime 修改文件过期时间，当前时间算起来，有效时间为 days 天
func (m *MediaAssetClient) ModifyExpireTime(mediaID uint64, days int) (requestID string, err error) {
	if m.Port == 0 {
		m.Port = 80
	}
	action := "ModifyExpireTime"
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
	req := &request.ModifyExpireTimeRequest{}
	req.TIBusinessID = m.TIBusinessID
	req.TIProjectID = m.TIProjectID
	req.MediaID = mediaID
	req.Days = int32(days)
	req.Action = action
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
	timeSleep := 50 * time.Millisecond
	rsp := &response.ModifyExpireTimeResponse{}
	for i := 0; i < maxTry; i++ {
		err = mediaassetservice.HttpPost(uri, header, req, rsp)
		if rsp.Response.ApiError != nil {
			bys, _ := json.Marshal(rsp)
			err = errors.New("Response error: " + string(bys))
		}
		if err == nil {
			break
		}
		time.Sleep(timeSleep)
		timeSleep *= 2
	}
	return rsp.Response.RequestID, err
}
