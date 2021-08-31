package mediaassetsdk

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"sync"

	"github.com/Tencent-Ti/ti-sign-go/tisign"
	"github.com/Tencent-media-asset-system-sdk/media-asset-go-sdk/common"
	mediaassetservice "github.com/Tencent-media-asset-system-sdk/media-asset-go-sdk/media_asset_service"
	"github.com/Tencent-media-asset-system-sdk/media-asset-go-sdk/media_asset_service/request"
	"github.com/Tencent-media-asset-system-sdk/media-asset-go-sdk/media_asset_service/response"
	ants "github.com/panjf2000/ants/v2"
)

// 上传分辨大小 10M
const BloackSzie = 32 * 1024 * 1024

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
	rsp := &response.ApplyUploadResponse{}
	for i := 0; i < maxTry; i++ {
		err = mediaassetservice.HttpPost(uri, header, req, rsp)
		if rsp.Response.ApiError != nil {
			bys, _ := json.Marshal(rsp)
			err = errors.New("Response error: " + string(bys))
		}
		if err == nil {
			break
		}
	}
	return rsp.Response.MediaID, rsp.Response.Bucket, rsp.Response.Key,
		rsp.Response.UploadId, rsp.Response.RequestID, err
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
	rsp := &response.CommitUploadResponse{}
	for i := 0; i < maxTry; i++ {
		err = mediaassetservice.HttpPost(uri, header, req, rsp)
		if rsp.Response.ApiError != nil {
			bys, _ := json.Marshal(rsp)
			err = errors.New("Response error: " + string(bys))
		}
		if err == nil {
			break
		}
		fmt.Println("Commit try ", i+1, " error: ", err.Error())
	}
	return rsp.Response.RequestID, err
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
		// canonicalQueryString = url.QueryEscape(canonicalQueryString)
		uri := ""
		header := map[string]string{}
		if m.Inner {
			uri = fmt.Sprintf("%s/UploadPart?%s", m.InnerFileManagerEndPoint, canonicalQueryString)
			header = nil
		} else {
			uri = fmt.Sprintf("http://%s:%d/FileManager/UploadPart?%s", m.Host, m.Port, canonicalQueryString)
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
			_, err = mediaassetservice.UploadPart(header, uri, filebuf)
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
	stat, err := os.Stat(filePath)
	if err == nil && stat.Size() <= BloackSzie {
		// 不需要分片
		filebuf, err := ioutil.ReadAll(file)
		if err != nil {
			return err
		}
		h := md5.New()
		h.Write(filebuf)
		md5sum := hex.EncodeToString(h.Sum(nil))
		canonicalQueryString := fmt.Sprintf("useJson=true&Bucket=%s&Key=%s&Content-MD5=%s", bucket, key, md5sum)
		uri := ""
		header := map[string]string{}
		if m.Inner {
			uri = fmt.Sprintf("%s/PutObject?%s", m.InnerFileManagerEndPoint, canonicalQueryString)
			header = nil
		} else {
			uri = fmt.Sprintf("http://%s:%d/FileManager/PutObject?%s", m.Host, m.Port, canonicalQueryString)
			headerContent := tisign.HttpHeaderContent{
				XTCAction:   "PutObject",               // 请求接口
				XTCService:  "app-cdn4aowk",             // 接口所属服务名
				XTCVersion:  "2021-02-26",               // 接口版本
				ContentType: "application/octet-stream", // http请求的content-type, 当前网关只支持: application/json  multipart/form-data
				HttpMethod:  "PUT",                      // http请求方法，当前网关只支持: POST GET
				Host:        m.Host,                     // 访问网关的host
			}
			ts := tisign.NewTiSign(headerContent, m.SecretID, m.SecretKey)
			header, _ = ts.CreateSignatureInfo()
			maxTry := 5
			for i := 0; i < maxTry; i++ {
				_, err = mediaassetservice.PutObject(header, uri, filebuf)
				if err == nil {
					break
				}
			}
			return err
		}
	} else {
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
			wg.Add(1)
			filebuf := make([]byte, n)
			copy(filebuf, buffer)
			_ = pool.Invoke([]interface{}{filebuf, partNumber})
			partNumber += 1
		}
		wg.Wait()
		if err != nil {
			return err
		}
	}
	return nil
}

func (m MediaAssetClient) doUploadBuf(buf []byte, key, bucket, uploadID string, coroutineNum int) (err error) {
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
		// canonicalQueryString = url.QueryEscape(canonicalQueryString)
		uri := ""
		header := map[string]string{}
		if m.Inner {
			uri = fmt.Sprintf("%s/UploadPart?%s", m.InnerFileManagerEndPoint, canonicalQueryString)
			header = nil
		} else {
			uri = fmt.Sprintf("http://%s:%d/FileManager/UploadPart?%s", m.Host, m.Port, canonicalQueryString)
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
			_, err = mediaassetservice.UploadPart(header, uri, filebuf)
			if err == nil {
				break
			}
		}
		wg.Done()
	})
	defer pool.Release()
	if len(buf) <= BloackSzie {
		fmt.Println("hhhhh")
		h := md5.New()
		h.Write(buf)
		md5sum := hex.EncodeToString(h.Sum(nil))
		canonicalQueryString := fmt.Sprintf("useJson=true&Bucket=%s&Key=%s&Content-MD5=%s", bucket, key, md5sum)
		uri := ""
		header := map[string]string{}
		if m.Inner {
			uri = fmt.Sprintf("%s/PutObject?%s", m.InnerFileManagerEndPoint, canonicalQueryString)
			header = nil
		} else {
			uri = fmt.Sprintf("http://%s:%d/FileManager/PutObject?%s", m.Host, m.Port, canonicalQueryString)
			headerContent := tisign.HttpHeaderContent{
				XTCAction:   "PutObject",               // 请求接口
				XTCService:  "app-cdn4aowk",             // 接口所属服务名
				XTCVersion:  "2021-02-26",               // 接口版本
				ContentType: "application/octet-stream", // http请求的content-type, 当前网关只支持: application/json  multipart/form-data
				HttpMethod:  "PUT",                      // http请求方法，当前网关只支持: POST GET
				Host:        m.Host,                     // 访问网关的host
			}
			ts := tisign.NewTiSign(headerContent, m.SecretID, m.SecretKey)
			header, _ = ts.CreateSignatureInfo()
			maxTry := 5
			for i := 0; i < maxTry; i++ {
				_, err = mediaassetservice.PutObject(header, uri, buf)
				if err == nil {
					break
				}
			}
			return err
		}
	} else {
		// Submit uploadpart one by one.
		partNumber := 1
		for i := 0; i < len(buf); i += BloackSzie {
			end := i + BloackSzie
			if end > len(buf) {
				end = len(buf)
			}
			wg.Add(1)
			_ = pool.Invoke([]interface{}{buf[i:end], partNumber})
			partNumber += 1
		}
		wg.Wait()
		if err != nil {
			return err
		}
	}
	return nil
}

// UploadFile 上传文件
// filePath 文件路径
// coroutineNum 上传最大并发协程数
// mediaInfo request.MediaMeta 媒体的类型和标签信息
func (m *MediaAssetClient) UploadFile(filePath, mediaName string, mediaMeta request.MediaMeta, coroutineNum int) (
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
	mediaID, bucket, key, uploadID, requestID, e := m.applyUplod(mediaName, mediaMeta, uint64(fileSize))
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

	// 第二步, 上传文件
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
	if mediaSet[0].Status != "上传完成" && mediaSet[0].Status != "验证素材中" {
		err = errors.New("素材错误, " + mediaSet[0].FailedReason)
	}
	return mediaSet[0], requestIDSet, err
}

// UploadBuf 上传内存文件
// filePath 文件路径
// coroutineNum 上传最大并发协程数
// mediaInfo request.MediaMeta 媒体的类型和标签信息
func (m *MediaAssetClient) UploadBuf(buf []byte, mediaName string, mediaMeta request.MediaMeta, coroutineNum int) (
	media *response.MediaInfo, requestIDSet []string, err error) {
	if m.Port == 0 {
		m.Port = 80
	}

	fileSize := len(buf)
	mediaID, bucket, key, uploadID, requestID, e := m.applyUplod(mediaName, mediaMeta, uint64(fileSize))
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

	// 第二步, 上传buf
	err = m.doUploadBuf(buf, key, bucket, uploadID, coroutineNum)
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
	if mediaSet[0].Status != "上传完成" && mediaSet[0].Status != "验证素材中" {
		err = errors.New("素材错误, " + mediaSet[0].FailedReason)
	}
	return mediaSet[0], requestIDSet, err
}
