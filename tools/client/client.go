package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

func DoClient(req *http.Request, timeOut time.Duration) (int, []byte, error) {
	client := &http.Client{
		Timeout: timeOut,
	}
	resp, err := client.Do(req)
	if err != nil {
		return 0, nil, err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return 0, nil, err
	}
	return resp.StatusCode, data, nil
}

func resultUnmarshal(data []byte, result interface{}) error {
	err := json.Unmarshal(data, result)
	if err != nil {
		return err
	}
	return nil
}

func PostNewNotify(method, notifyURL string, query map[string]string, bodyByte []byte) ([]byte, error) {

	SetQuery(&notifyURL, query)

	client := &http.Client{
		Timeout: time.Duration(5) * time.Minute,
	}

	req, err := http.NewRequest(method, notifyURL, bytes.NewBuffer(bodyByte))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	var result []byte
	defer resp.Body.Close()
	result, err = ioutil.ReadAll(resp.Body)
	if resp.StatusCode/200 != 1 {
		return nil, fmt.Errorf("err status code %v", resp.StatusCode)
	}
	return result, nil
}

func SetQuery(URL *string, kv map[string]string) {
	newURL, err := url.Parse(*URL)
	if err != nil {
		return
	}

	vs := newURL.Query()
	for k, v := range kv {
		vs.Set(k, v)
	}

	newURL.RawQuery = vs.Encode()

	*URL = newURL.String()
	return
}
