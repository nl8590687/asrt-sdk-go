package common

var (
	APIStatusCodeOK                 int = 200000 // OK
	APIStatusCodePartOK             int = 206000 // 部分识别结果
	APIStatusCodeClientError        int = 400000
	APIStatusCodeClientErrorFormat  int = 400001 // 请求数据格式错误
	APIStatusCodeClientErrorConfig  int = 400002 // 请求数据配置不支持
	APIStatusCodeServerError        int = 500000
	APIStatusCodeServerErrorRunning int = 500001 // 服务器运行中出错
)

// AsrtAPIResponse ASRT语音识别API响应类
type AsrtAPIResponse struct {
	StatusCode    int         `json:"status_code"`
	StatucMesaage string      `json:"status_message"`
	Result        interface{} `json:"result"`
}

// AsrtAPISpeechRequest ASRT语音识别API语音数据请求类
type AsrtAPISpeechRequest struct {
	Samples    string `json:"samples"`
	SampleRate int    `json:"sample_rate"`
	Channels   int    `json:"channels"`
	ByteWidth  int    `json:"byte_width"`
}

// AsrtAPILanguageRequest ASRT语音识别API语言模型请求类
type AsrtAPILanguageRequest struct {
	SequencePinyin []string `json:"sequence_pinyin"`
}
