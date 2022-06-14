package main

import (
	"fmt"
	"time"

	"github.com/nl8590687/asrt-sdk-go/common"
	"github.com/nl8590687/asrt-sdk-go/sdk"
)

func main() {
	httpDemo()
	grpcDemo()
}

func httpDemo() {
	// 初始化
	host := "127.0.0.1"
	port := "20001"
	protocol := "http"

	sr := sdk.GetSpeechRecognizer(host, port, protocol)
	// ======================================================
	// 识别文件
	filename := "testData/data1.wav"
	resultFile, err := sr.RecogniteFile(filename)
	if err != nil {
		fmt.Println(err)
	}

	for index, res := range resultFile {
		fmt.Println("Wav文件语音识别结果 ", index, ":", res.Result)
	}

	byteData := sdk.LoadFile(filename)
	wave, err := sdk.DecodeWav(byteData)
	if err != nil {
		fmt.Println(err)
	}
	// ======================================================
	// 识别一段Wave音频序列
	result, err := sr.Recognite(wave.GetRawSamples(), wave.FrameRate, wave.Channels, wave.SampleWidth)
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

func grpcDemo() {
	// 初始化
	host := "127.0.0.1"
	port := "20002"
	protocol := "grpc"

	sr := sdk.GetSpeechRecognizer(host, port, protocol)
	fmt.Println("sr:", sr)
	// ======================================================
	// 识别文件
	filename := "testData/data1.wav"
	resultFile, err := sr.RecogniteFile(filename)
	if err != nil {
		fmt.Println(err)
	}

	for index, res := range resultFile {
		fmt.Println("Wav文件语音识别结果 ", index, ":", res.Result)
	}

	byteData := sdk.LoadFile(filename)
	wave, err := sdk.DecodeWav(byteData)
	if err != nil {
		fmt.Println(err)
	}
	// ======================================================
	// 识别一段Wave音频序列
	result, err := sr.Recognite(wave.GetRawSamples(), wave.FrameRate, wave.Channels, wave.SampleWidth)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("语音识别结果：", result.Result)
	// ======================================================
	// 识别一段长Wave音频序列
	longSample := wave.GetRawSamples()
	longSample = append(longSample, wave.GetRawSamples()...)
	resultLong, err := sr.RecogniteLong(longSample, wave.FrameRate, wave.Channels, wave.SampleWidth)
	if err != nil {
		fmt.Println(err)
	}

	for index, res := range resultLong {
		fmt.Println("长文件语音识别结果 ", index, ":", res.Result)
	}
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
	for i := 0; i < len(result.Result.([]string)); i += 1 {
		pinyinResult = append(pinyinResult, result.Result.([]string)[i])
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
	// ======================================================
	// 调用ASRT grpc接口流式识别
	var wavChannel = make(chan *common.Wav, 5)
	var recognitionResult = make(chan *common.AsrtAPIResponse, 5)
	sendFunction := func() {
		var index int
		for index = 0; index < 10; index += 1 {
			time.Sleep(2 * time.Second)
			wavChannel <- wave
		}
		close(wavChannel)
	}
	go sendFunction()
	var asrResult string
	var tmpAsrResult string
	recvFunction := func() {
		for value := range recognitionResult {
			fmt.Println("流式解码结果：", value.StatusCode, value.Result, value.StatucMesaage)
			if value.StatusCode == common.APIStatusCodeOK {
				tmpAsrResult = ""
				asrResult += value.Result.(string)
			} else if value.StatusCode == common.APIStatusCodePartOK {
				tmpAsrResult = value.Result.(string)
			}
			fmt.Println("语音识别文本：", asrResult+tmpAsrResult)
		}
	}
	go recvFunction()
	err = sr.(*sdk.GRPCSpeechRecognizer).RecogniteStream(wavChannel, recognitionResult)
	fmt.Println("流式识别完毕")
	if err != nil {
		fmt.Println(err)
	}
	close(recognitionResult)
}
