package media_asset_service

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/Tencent-media-asset-system-sdk/media-asset-go-sdk/media_asset_service/response"
)

// UploadPart 上传分片
func UploadPart(header map[string]string, uri string, filebuf []byte) (*response.UploadPartResponse, error) {

	rsp := &response.UploadPartResponse{}
	req, err := http.NewRequest(http.MethodPut, uri, bytes.NewBuffer(filebuf))
	if err != nil {
		errstr := fmt.Sprintf("[%s]UploadPart make NewRequest error: %s", uri, err.Error())
		return nil, errors.New(errstr)
	}
	if header == nil {
		req.Header.Set("Content-Type", "application/octet-stream")
	} else {
		for head, value := range header {
			req.Header.Set(head, value)
		}
	}
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		errstr := fmt.Sprintf("[%s]UploadPart do request error: %s", uri, err.Error())
		return nil, errors.New(errstr)
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		errstr := fmt.Sprintf("[%s]UploadPart status error, statuscode: %d, reason: %s",
			uri, response.StatusCode, response.Status)
		return nil, errors.New(errstr)
	}
	data, _ := ioutil.ReadAll(response.Body)
	err = json.Unmarshal(data, rsp)
	if err != nil {
		errstr := fmt.Sprintf("[%s]UploadPart response protocol error %s", uri, err.Error())
		return nil, errors.New(errstr)
	}
	return rsp, nil
}
