package sdk

import (
	"github.com/nl8590687/asrt-sdk-go/common"
)

type ISpeechRecognizer interface {
	Recognite(wavData []byte, frameRate int, channels int, byteWidth int) (*common.AsrtApiResponse, error)
	RecogniteSpeech(wavData []byte, frameRate int, channels int, byteWidth int) (*common.AsrtApiResponse, error)
	RecogniteLanguage(sequencePinyin []string) (*common.AsrtApiResponse, error)
	RecogniteFile(filename string) (*common.AsrtApiResponse, error)
}

type BaseSpeechRecognizer struct {
	Host     string
	Port     string
	Protocol string
}
