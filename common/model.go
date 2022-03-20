package common

var (
	API_STATUS_CODE_OK                   int = 200000 // OK
	API_STATUS_CODE_CLIENT_ERROR         int = 400000
	API_STATUS_CODE_CLIENT_ERROR_FORMAT  int = 400001 // 请求数据格式错误
	API_STATUS_CODE_CLIENT_ERROR_CONFIG  int = 400002 // 请求数据配置不支持
	API_STATUS_CODE_SERVER_ERROR         int = 500000
	API_STATUS_CODE_SERVER_ERROR_RUNNING int = 500001 // 服务器运行中出错
)

type AsrtApiResponse struct {
	StatusCode    int         `json:"status_code"`
	StatucMesaage string      `json:"status_message"`
	Result        interface{} `json:"result"`
}

type AsrtApiSpeechRequest struct {
	Samples    string `json:"samples"`
	SampleRate int    `json:"sample_rate"`
	Channels   int    `json:"channels"`
	ByteWidth  int    `json:"byte_width"`
}

type AsrtApiLanguageRequest struct {
	SequencePinyin []string `json:"sequence_pinyin"`
}
