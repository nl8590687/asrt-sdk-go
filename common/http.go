package common

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

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

func UrlEncode(text string) string {
	var urlStr string = text
	escapeUrl := url.QueryEscape(urlStr)
	return escapeUrl
}

func UrlDecode(text string) string {
	enEscapeUrl, err := url.QueryUnescape(text)
	if err != nil {
		log.Fatalln("error: URL decode failed.", err)
		return ""
	}

	return enEscapeUrl
}
