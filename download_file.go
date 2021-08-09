package mediaassetsdk

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
	"strconv"

	"github.com/Tencent-Ti/ti-sign-go/tisign"
	"github.com/tidwall/gjson"
)

// DownloadFile 通过媒体信息返回的url下载文件到本地
func (m *MediaAssetClient) DownloadFile(downloadURL, dir, fileName string) (err error) {
	if m.Port == 0 {
		m.Port = 80
	}
	if m.Inner {
		// 内网调用直接软链，挂载公共数据盘
		u, err := url.Parse(downloadURL)
		if err != nil {
			return errors.New("DownloadFileInner failed, parse url failed, error: " + err.Error() + " url: " + downloadURL)
		}
		filePath := u.Query().Get("Key")
		filePath = path.Join(m.InnerDataDir, filePath)
		if filePath == "" {
			return errors.New("DownloadFileInner failed, filePath is null, url: " + downloadURL)
		}
		if _, err := os.Stat(filePath); err != nil {
			return errors.New("DownloadFileInner failed, file " + filePath + " stat error : " + err.Error())
		}
		if err := os.MkdirAll(dir, 0766); err != nil {
			return errors.New("DownloadFileInner " + downloadURL + " failed! MkdirAll error: " + err.Error())
		}
		if err := os.Symlink(filePath, path.Join(dir, fileName)); err != nil {
			return errors.New("DownloadFileInner failed, Symlink failed, error: " + err.Error())
		}
		return nil
	} else {
		uri := fmt.Sprintf("http://%s:%d%s", m.Host, m.Port, downloadURL)
		action := "DownloadFile"
		service := "app-cdn4aowk"
		version := "2021-02-26"
		headerContent := tisign.HttpHeaderContent{
			XTCAction:   action,             // 请求接口
			XTCService:  service,            // 接口所属服务名
			XTCVersion:  version,            // 接口版本
			ContentType: "application/json", // http请求的content-type, 当前网关只支持: application/json  multipart/form-data
			HttpMethod:  "GET",              // http请求方法，当前网关只支持: POST GET
			Host:        m.Host,             // 访问网关的host
		}
		ts := tisign.NewTiSign(headerContent, m.SecretID, m.SecretKey)
		header, _ := ts.CreateSignatureInfo()
		req, err := http.NewRequest("GET", uri, nil)
		if err != nil {
			return errors.New("DownloadFile " + uri + " make newrequest failed! " + err.Error())
		}
		for head, value := range header {
			req.Header.Set(head, value)
		}
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			return errors.New("DownloadFile " + uri + " do request failed! " + err.Error())
		}
		if res.StatusCode != 200 {
			return errors.New("DownloadFile " + uri + " failed! " + res.Status)
		}
		if err := os.MkdirAll(dir, 0766); err != nil {
			return errors.New("DownloadFile " + uri + " failed! " + err.Error())
		}
		filePath := dir + "/" + fileName
		file, err := os.Create(filePath)
		if err != nil {
			return errors.New("DownloadFile " + uri + " create file failed! " + err.Error())
		}
		if length, err := io.Copy(file, res.Body); err != nil {
			os.Remove(filePath)
			return errors.New("DownloadFile " + uri + " copy failed! " + err.Error() +
				" copied length: " + strconv.Itoa(int(length)))
		}
		file.Close()
		_, err = os.Stat(filePath)
		if err != nil {
			os.Remove(filePath)
			return errors.New("DownloadFile " + uri + " failed! " + err.Error())
		}
		return nil
	}
}

// DownloadToBuf 通过媒体信息返回的url下载文件到内存
func (m *MediaAssetClient) DownloadToBuf(downloadURL string) (buf []byte, err error) {
	if m.Port == 0 {
		m.Port = 80
	}
	action := "DownloadFile"
	service := "app-cdn4aowk"
	version := "2021-02-26"
	headerContent := tisign.HttpHeaderContent{
		XTCAction:   action,                     // 请求接口
		XTCService:  service,                    // 接口所属服务名
		XTCVersion:  version,                    // 接口版本
		ContentType: "application/octet-stream", // http请求的content-type, 当前网关只支持: application/json  multipart/form-data
		HttpMethod:  "GET",                      // http请求方法，当前网关只支持: POST GET
		Host:        m.Host,                     // 访问网关的host
	}
	uri := ""
	header := map[string]string{}
	if m.Inner {
		uri = fmt.Sprintf("%s%s", m.InnerFileStaticEndPoint, downloadURL)
		header["Content-Type"] = "application/octet-stream"
		header["X-TC-Uin"] = m.InnerUserName
	} else {
		uri = fmt.Sprintf("http://%s:%d%s", m.Host, m.Port, downloadURL)
		ts := tisign.NewTiSign(headerContent, m.SecretID, m.SecretKey)
		header, _ = ts.CreateSignatureInfo()
	}
	maxTry := 5
	fmt.Println(uri)
	for i := 0; i < maxTry; i++ {
		req, e := http.NewRequest("GET", uri, nil)
		if e != nil {
			err = errors.New("DownloadToBuf " + uri + " make newrequest failed! " + e.Error())
			continue
		}
		for head, value := range header {
			req.Header.Set(head, value)
		}
		res, e := http.DefaultClient.Do(req)
		if e != nil {
			err = errors.New("DownloadToBuf " + uri + " do request failed! " + e.Error())
			continue
		}
		if res.StatusCode != 200 {
			err = errors.New("DownloadToBuf " + uri + " failed! " + res.Status)
			continue
		}
		defer res.Body.Close()
		buf, err = ioutil.ReadAll(res.Body)
		data := gjson.ParseBytes(buf)
		if data.Get("Response.Error").Exists() {
			err = errors.New("DownloadToBuf " + uri + " failed! response error, data: " + string(buf))
		}
		if err == nil && len(buf) > 0 {
			break
		}
	}
	if err != nil {
		return nil, err
	}
	// 防止 http 请求返回空数据
	if len(buf) == 0 {
		return nil, errors.New("http request return null file")
	}
	return buf, nil
}
