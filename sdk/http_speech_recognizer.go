package sdk

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/nl8590687/asrt-sdk-go/common"
)

// HttpSpeechRecognizer 调用ASRT语音识别系统HTTP+JSON协议接口的语音识别类
type HttpSpeechRecognizer struct {
	BaseSpeechRecognizer
	// SubPath HTTP协议资源子路径，默认为""
	SubPath string
}

// NewHttpSpeechRecognizer 构造一个用于调用http+json协议接口的语音识别类实例对象
func NewHttpSpeechRecognizer(host string, port string, protocol string) *HttpSpeechRecognizer {
	base := BaseSpeechRecognizer{
		Host:     host,
		Port:     port,
		Protocol: protocol,
	}
	httpSpeechRecognizer := HttpSpeechRecognizer{
		BaseSpeechRecognizer: base,
		SubPath:              "",
	}

	return &httpSpeechRecognizer
}

func (h *HttpSpeechRecognizer) getUrl() string {
	return fmt.Sprintf("%s://%s:%s%s", h.Protocol, h.Host, h.Port, h.SubPath)
}

// Recognite 调用ASRT语音识别
func (h *HttpSpeechRecognizer) Recognite(wavData []byte, frameRate int, channels int, byteWidth int) (*common.AsrtApiResponse, error) {
	requestBody := common.AsrtApiSpeechRequest{
		Samples:    common.BytesToBase64(wavData),
		SampleRate: frameRate,
		Channels:   channels,
		ByteWidth:  byteWidth,
	}

	byteForm, err := json.Marshal(requestBody)

	if err != nil {
		return nil, err
	}

	contentType := "application/json"
	url := fmt.Sprintf("%s/all", h.getUrl())
	rspBody, err := common.SendHttpRequestPost(url, byteForm, contentType)
	if err != nil {
		return nil, err
	}

	responseBody := common.AsrtApiResponse{}
	err = json.Unmarshal(rspBody, &responseBody)
	if err != nil {
		return nil, err
	}

	log.Println("info: recv: ", responseBody)
	return &responseBody, nil
}

// RecogniteSpeech 调用ASRT语音识别声学模型
func (h *HttpSpeechRecognizer) RecogniteSpeech(wavData []byte, frameRate int, channels int, byteWidth int) (*common.AsrtApiResponse, error) {
	requestBody := common.AsrtApiSpeechRequest{
		Samples:    common.BytesToBase64(wavData),
		SampleRate: frameRate,
		Channels:   channels,
		ByteWidth:  byteWidth,
	}

	byteForm, err := json.Marshal(requestBody)
	if err != nil {
		return nil, err
	}

	contentType := "application/json"
	url := fmt.Sprintf("%s/speech", h.getUrl())
	rspBody, err := common.SendHttpRequestPost(url, byteForm, contentType)
	if err != nil {
		return nil, err
	}

	responseBody := common.AsrtApiResponse{}
	err = json.Unmarshal(rspBody, &responseBody)
	if err != nil {
		return nil, err
	}

	log.Println("info: recv: ", responseBody)
	return &responseBody, nil
}

// RecogniteLanguage 调用ASRT语音识别语言模型
func (h *HttpSpeechRecognizer) RecogniteLanguage(sequencePinyin []string) (*common.AsrtApiResponse, error) {
	requestBody := common.AsrtApiLanguageRequest{
		SequencePinyin: sequencePinyin,
	}

	byteForm, err := json.Marshal(requestBody)
	if err != nil {
		return nil, err
	}

	contentType := "application/json"
	url := fmt.Sprintf("%s/language", h.getUrl())
	rspBody, err := common.SendHttpRequestPost(url, byteForm, contentType)
	if err != nil {
		return nil, err
	}

	responseBody := common.AsrtApiResponse{}
	err = json.Unmarshal(rspBody, &responseBody)
	if err != nil {
		return nil, err
	}

	log.Println("info: recv: ", responseBody)
	return &responseBody, nil
}

// RecogniteFile 调用ASRT语音识别来识别指定文件名的音频文件
func (h *HttpSpeechRecognizer) RecogniteFile(filename string) (*common.AsrtApiResponse, error) {
	binData := common.ReadBinFile(filename)
	wavAudio := common.Wav{}
	err := wavAudio.Deserialize(binData)
	if err != nil {
		return nil, err
	}

	byteData := wavAudio.GetRawSamples()
	rsp, err := h.Recognite(byteData, wavAudio.FrameRate, wavAudio.Channels, wavAudio.SampleWidth)

	return rsp, err
}
