package sdk

import (
	"github.com/nl8590687/asrt-sdk-go/common"
)

// ISpeechRecognizer ASRT语音识别SDK语音识别抽象接口
type ISpeechRecognizer interface {
	// Recognite 调用ASRT语音识别
	Recognite(wavData []byte, frameRate int, channels int, byteWidth int) (*common.AsrtAPIResponse, error)
	// RecogniteSpeech 调用ASRT语音识别声学模型
	RecogniteSpeech(wavData []byte, frameRate int, channels int, byteWidth int) (*common.AsrtAPIResponse, error)
	// RecogniteLanguage 调用ASRT语音识别语言模型
	RecogniteLanguage(sequencePinyin []string) (*common.AsrtAPIResponse, error)
	// RecogniteLong
	RecogniteLong(wavData []byte, frameRate int, channels int, byteWidth int) ([]*common.AsrtAPIResponse, error)
	// RecogniteFile 调用ASRT语音识别来识别指定文件名的音频文件
	RecogniteFile(filename string) ([]*common.AsrtAPIResponse, error)
}

// BaseSpeechRecognizer ASRT语音识别SDK语音识别基类
type BaseSpeechRecognizer struct {
	// Host 主机域名或IP
	Host string
	// Port 主机端口号
	Port string
	// Protocol 网络协议
	Protocol string
}
