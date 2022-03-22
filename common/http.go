package common

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

// SendHttpRequestGet 发送HTTP GET请求
func SendHttpRequestGet(url string) ([]byte, error) {
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

// SendHttpRequestPost 发送HTTP POST请求
func SendHttpRequestPost(url string, bytesForm []byte, contentType string) ([]byte, error) {
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

// URLEncode URL编码
func URLEncode(text string) string {
	var urlStr = text
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
