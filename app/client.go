package app

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func DoRequestDecode(client *http.Client, req *http.Request, v any) (resp *http.Response, body []byte, err error) {
	resp, body, err = DoRequestGetBody(client, req)
	if err == nil {
		err = json.Unmarshal(body, v)
	}
	return
}

func DoRequestGetBody(client *http.Client, req *http.Request) (resp *http.Response, body []byte, err error) {
	resp, err = client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	var buf bytes.Buffer
	io.Copy(&buf, resp.Body)
	body = buf.Bytes()
	if resp.StatusCode < 200 || resp.StatusCode > 300 {
		err = fmt.Errorf("error code %d", resp.StatusCode)
	}
	return
}
