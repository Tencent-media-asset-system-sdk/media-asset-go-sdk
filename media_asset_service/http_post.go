package mediaassetservice

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

// HttpPost http post 请求
func HttpPost(uri string, header map[string]string, req interface{}, rsp interface{}) (err error) {
	body, _ := json.Marshal(req)
	reqCtx, err := http.NewRequest(http.MethodPost, uri, bytes.NewReader(body))
	if err != nil {
		errstr := fmt.Sprintf("HttpPost make NewRequest error: %s", err.Error())
		return errors.New(errstr)
	}
	if header == nil {
		reqCtx.Header.Set("Content-Type", "application/json")
	} else {
		for head, value := range header {
			reqCtx.Header.Set(head, value)
		}
	}
	response, err := http.DefaultClient.Do(reqCtx)
	if err != nil {
		errstr := fmt.Sprintf("HttpPost do request error: %s", err.Error())
		return errors.New(errstr)
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		errstr := fmt.Sprintf("HttpPost response status error, statuscode: %d, reason: %s",
			response.StatusCode, response.Status)
		return errors.New(errstr)
	}
	data, _ := ioutil.ReadAll(response.Body)
	fixedData := fmt.Sprintf("{\"Response\":%s}", string(data))
	err = json.Unmarshal([]byte(fixedData), rsp)
	if err != nil {
		fmt.Print(err)
		errstr := fmt.Sprintf("HttpPost response protocol error, data: %s", fixedData)
		return errors.New(errstr)
	}
	return nil
}
