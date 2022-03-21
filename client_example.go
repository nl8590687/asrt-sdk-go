package main

import (
	"fmt"

	"github.com/nl8590687/asrt-sdk-go/sdk"
)

func main() {
	// 初始化
	host := "127.0.0.1"
	port := "20001"
	protocol := "http"

	sr := sdk.GetSpeechRecognizer(host, port, protocol)
	// ======================================================
	// 识别文件
	filename := "testData/data1.wav"
	result, err := sr.RecogniteFile(filename)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("语音识别结果：", result.Result)

	byteData := sdk.LoadFile(filename)
	wave, err := sdk.DecodeWav(byteData)
	if err != nil {
		fmt.Println(err)
	}
	// ======================================================
	// 识别一段Wave音频序列
	result, err = sr.Recognite(wave.GetRawSamples(), wave.FrameRate, wave.Channels, wave.SampleWidth)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("语音识别结果：", result.Result)
	// ======================================================
	// 调用声学模型识别一段Wave音频序列
	result, err = sr.RecogniteSpeech(wave.GetRawSamples(), wave.FrameRate, wave.Channels, wave.SampleWidth)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("语音识别声学模型结果：", result.Result)
	// ======================================================
	// 调用语言模型1
	pinyinResult := []string{}
	for i := 0; i < len(result.Result.([]interface{})); i += 1 {
		pinyinResult = append(pinyinResult, result.Result.([]interface{})[i].(string))
	}

	result, err = sr.RecogniteLanguage(pinyinResult)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("语言模型结果：", result.Result)
	// ======================================================
	// 调用语言模型2
	sequencePinyin := []string{"ni3", "hao3", "a1"}
	result, err = sr.RecogniteLanguage(sequencePinyin)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("语言模型结果：", result.Result)
}
