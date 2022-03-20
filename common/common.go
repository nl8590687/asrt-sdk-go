// package common 一些有用的功能库
package common

import (
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"os"
)

// Base64ToBytes base64编码数据转换为[]byte字节数组
func Base64ToBytes(base64Data string) []byte {
	data, err := base64.StdEncoding.DecodeString(base64Data)
	if err != nil {
		log.Println(err)
	}

	return data
}

// Base64ToBytes_UrlSafe URL安全的base64编码数据转换为[]byte字节数组
func Base64ToBytes_UrlSafe(base64Data string) []byte {
	data, err := base64.URLEncoding.DecodeString(base64Data)
	if err != nil {
		log.Println(err)
	}

	return data
}

// BytesToBase64 []byte字节数组转换为base64编码数据
func BytesToBase64(bytesData []byte) string {
	encodedMsg := base64.StdEncoding.EncodeToString(bytesData)
	return encodedMsg
}

// BytesToBase64_UrlSafe []byte字节数组转换为URL安全的base64编码数据
func BytesToBase64_UrlSafe(bytesData []byte) string {
	encodedMsg := base64.URLEncoding.EncodeToString(bytesData)
	return encodedMsg
}

func readBinFile(filename string) []byte {
	fp, err := os.Open(filename)
	if err != nil {
		log.Println(err)
	}
	defer fp.Close()

	buff := make([]byte, 1024000*64) // 文件长度
	var length int

	for {
		lens, err := fp.Read(buff)
		length += lens
		if err == io.EOF || lens < 0 {
			break
		}
	}

	return buff[0:length]
}

func readTextFile(filename string) string {
	return string(readBinFile(filename))
}

func writeBinFile(filename string, data []byte) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("error: 文件创建失败, %s", err.Error())
	}
	defer file.Close()

	_, err = file.Write(data)
	if err != nil {
		return fmt.Errorf("error: 写入文件失败, %s", err.Error())
	}

	return nil
}

func SaveWaveObject(filename string, wave Wav) error {
	wavBytes, err := wave.Serialize()
	if err != nil {
		return err
	}

	err = writeBinFile(filename, wavBytes)
	if err != nil {
		return err
	}

	return nil
}

func ReadBinFile(filename string) []byte {
	return readBinFile(filename)
}
