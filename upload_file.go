package media_asset_sdk

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/url"
	"os"
	"strconv"
	"sync"

	"github.com/Tencent-Ti/ti-sign-go/tisign"
	"github.com/Tencent-media-asset-system-sdk/media-asset-go-sdk/common/utils"
	"github.com/Tencent-media-asset-system-sdk/media-asset-go-sdk/media_asset_service"
	"github.com/Tencent-media-asset-system-sdk/media-asset-go-sdk/media_asset_service/request"
	"github.com/Tencent-media-asset-system-sdk/media-asset-go-sdk/media_asset_service/response"
	ants "github.com/panjf2000/ants/v2"
)

// 上传分辨大小 10M
const BloackSzie = 10 * 1024 * 1024

func (m MediaAssetClient) applyUplod(mediaName string, mediaMeta request.MediaMeta, fileSize uint64) (
	mediaID uint64, bucket, key, uploadId, requestID string, err error) {

	action := "ApplyUpload"
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
	req := &request.ApplyUploadRequest{}
	req.TIBusinessID = uint32(m.TIBusinessID)
	req.TIProjectID = uint32(m.TIProjectID)
	req.Name = mediaName
	req.MediaMeta = mediaMeta
	req.Size = strconv.FormatUint(fileSize, 10)
	req.Inner = m.Inner
	if m.Inner {
		req.RequestID = utils.GenerateRandomString(32)
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
	rsp := &response.ApplyUploadResponse{}
	for i := 0; i < maxTry; i++ {
		err = media_asset_service.HttpPost(uri, header, req, rsp)
		if rsp.ApiError != nil {
			bys, _ := json.Marshal(rsp)
			err = errors.New("Response error: " + string(bys))
		}
		if err == nil {
			break
		}
	}
	return rsp.MediaID, rsp.Bucket, rsp.Key, rsp.UploadId, rsp.RequestID, err
}

func (m MediaAssetClient) commitUpload(mediaID uint64, bucket, key, uploadID string) (
	requestID string, err error) {

	action := "CommitUpload"
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
	req := &request.CommitUploadRequest{}
	req.TIBusinessID = uint32(m.TIBusinessID)
	req.TIProjectID = uint32(m.TIProjectID)
	req.MediaID = mediaID
	req.Bucket = bucket
	req.Key = key
	req.UploadId = uploadID
	if m.Inner {
		req.RequestID = utils.GenerateRandomString(32)
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
	rsp := &response.CommitUploadResponse{}
	for i := 0; i < maxTry; i++ {
		err = media_asset_service.HttpPost(uri, header, req, rsp)
		if rsp.ApiError != nil {
			bys, _ := json.Marshal(rsp)
			err = errors.New("Response error: " + string(bys))
		}
		if err == nil {
			break
		}
	}
	return rsp.RequestID, err
}

func (m MediaAssetClient) doUpload(filePath, key, bucket, uploadID string, coroutineNum int) (err error) {
	defer ants.Release()
	wg := sync.WaitGroup{}
	// 线程池
	pool, _ := ants.NewPoolWithFunc(coroutineNum, func(v interface{}) {
		filebuf := v.([]interface{})[0].([]byte)
		partNumber := v.([]interface{})[1].(int)
		h := md5.New()
		h.Write(filebuf)
		md5sum := hex.EncodeToString(h.Sum(nil))
		canonicalQueryString := fmt.Sprintf("useJson=true&Bucket=%s&Key=%s&uploadId=%s&partNumber=%d&Content-MD5=%s",
			bucket, key, uploadID, partNumber, md5sum)
		canonicalQueryString = url.QueryEscape(canonicalQueryString)
		uri := fmt.Sprintf("http://%s:%d/UploadPart?%s", m.Host, m.Port, canonicalQueryString)
		header := map[string]string{}
		if m.Inner {
			header = nil
		} else {
			headerContent := tisign.HttpHeaderContent{
				XTCAction:   "UploadPart",               // 请求接口
				XTCService:  "app-cdn4aowk",             // 接口所属服务名
				XTCVersion:  "2021-02-26",               // 接口版本
				ContentType: "application/octet-stream", // http请求的content-type, 当前网关只支持: application/json  multipart/form-data
				HttpMethod:  "PUT",                      // http请求方法，当前网关只支持: POST GET
				Host:        m.Host,                     // 访问网关的host
			}
			ts := tisign.NewTiSign(headerContent, m.SecretID, m.SecretKey)
			header, _ = ts.CreateSignatureInfo()
		}
		maxTry := 5
		for i := 0; i < maxTry; i++ {
			_, err = media_asset_service.UploadPart(header, uri, filebuf)
			if err == nil {
				break
			}
		}
		wg.Done()
	})
	defer pool.Release()

	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	buffer := make([]byte, BloackSzie)
	// Submit uploadpart one by one.
	partNumber := 1
	for {
		n, err := file.Read(buffer)
		if err != nil && err != io.EOF {
			return err
		}
		if n == 0 {
			break
		}
		// 如下代码打印出每次读取的文件块(字节数)
		fmt.Println(string(buffer[:n]))
		wg.Add(1)
		filebuf := make([]byte, n)
		copy(filebuf, buffer)
		_ = pool.Invoke([]interface{}{filebuf, partNumber})
		partNumber += 1
	}
	wg.Wait()
	fmt.Println("Uploadfile success")
	return nil
}

// UploadFile 上传文件
// filePath 文件路径
// coroutineNum 上传最大并发协程数
// mediaInfo request.MediaMeta 媒体的类型和标签信息
func (m MediaAssetClient) UploadFile(filePath, mediaName string, MediaMeta request.MediaMeta, coroutineNum int) (
	media *response.MediaInfo, requestIDSet []string, err error) {
	if m.Port == 0 {
		m.Port = 80
	}
	// 第一步, 检查文件
	stat, e := os.Stat(filePath)
	if e != nil {
		return media, requestIDSet, e
	}
	fileSize := stat.Size()
	mediaID, key, bucket, uploadID, requestID, e := m.applyUplod(mediaName, media.MediaMeta, uint64(fileSize))
	if requestID != "" {
		requestIDSet = append(requestIDSet, requestID)
	}
	if e != nil {
		err = errors.New("UploadFile error in ApplyUpload: " + e.Error())
		return media, requestIDSet, err
	}
	defer func() {
		if err != nil {
			m.RemoveMedias([]uint64{mediaID})
		}
	}()

	// 第二步, 上传分片
	err = m.doUpload(filePath, key, bucket, uploadID, coroutineNum)
	if err != nil {
		err = errors.New("UploadFile error in UploadPart: " + err.Error())
		return media, requestIDSet, err
	}

	// 第三步, 确认上传
	reqID, e := m.commitUpload(mediaID, bucket, key, uploadID)
	if reqID != "" {
		requestIDSet = append(requestIDSet, reqID)
	}
	if e != nil {
		err = errors.New("UploadFile error in CommitUpload: " + e.Error())
		return media, requestIDSet, err
	}

	// 第四步，查询媒体信息
	mediaSet, reqID, e := m.DescribeMediaDetails([]uint64{mediaID})
	if reqID != "" {
		requestIDSet = append(requestIDSet, reqID)
	}
	if e != nil {
		err = errors.New("UploadFile error in DescribeMediaDetails: " + e.Error())
		return media, requestIDSet, err
	}
	if len(mediaSet) != 1 {
		err = errors.New("UploadFile error, DescribeMediaDetails return null mediaiInfo")
		return media, requestIDSet, err
	}
	return mediaSet[0], requestIDSet, err
}
