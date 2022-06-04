package common

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"runtime"
	"time"
)

var httpUserAgent string = fmt.Sprintf("%s%s%s%s%s", "ASRT-SDK client/", "v1",
	" (", runtime.Version(), ") (https://asrt.ailemon.net/)")

// SendHTTPRequestGet 发送HTTP GET请求
func SendHTTPRequestGet(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalln("error:", err)
		return nil, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln("error:", err)
		return nil, err
	}

	if resp.StatusCode == 200 {
		log.Println("info: http status 200 ok")
	} else {
		log.Println("warning: http status ", resp.StatusCode)
	}

	return body, nil
}

// SendHTTPRequestPost 发送HTTP POST请求
func SendHTTPRequestPost(url string, bytesForm []byte, contentType string) ([]byte, error) {
	bodyReader := bytes.NewReader(bytesForm)
	resp, err := http.Post(url, contentType, bodyReader)
	if err != nil {
		log.Fatalln("error:", err)
		return nil, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln("error:", err)
		return nil, err
	}

	if resp.StatusCode == 200 {
		log.Println("info: http status 200 ok")
	} else {
		log.Println("warning: http status ", resp.StatusCode)
	}

	return body, nil
}

// SendHTTPRequest 发送HTTP请求
func SendHTTPRequest(url string, method string,
	bytesBody []byte, contentType string,
) ([]byte, error) {
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	var bodyReader io.Reader = nil
	if method == "POST" {
		bodyReader = bytes.NewReader(bytesBody)
	}

	req, err := http.NewRequest(method, url, bodyReader)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", httpUserAgent)
	if len(contentType) > 0 {
		req.Header.Set("Content-Type", contentType)
	}

	rsp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer rsp.Body.Close()
	rspBody, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return nil, err
	}

	if rsp.StatusCode != 200 {
		log.Println("warning: http status ", rsp.StatusCode)
	}

	return rspBody, nil
}

// URLEncode URL编码
func URLEncode(text string) string {
	urlStr := text
	escapeURL := url.QueryEscape(urlStr)
	return escapeURL
}

// URLDecode URL解码
func URLDecode(text string) string {
	enEscapeURL, err := url.QueryUnescape(text)
	if err != nil {
		log.Fatalln("error: URL decode failed.", err)
		return ""
	}

	return enEscapeURL
}
