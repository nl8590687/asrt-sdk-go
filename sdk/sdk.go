package sdk

func GetSpeechRecognizer(host string, port string, protocol string) ISpeechRecognizer {
	if protocol == "http" || protocol == "https" {
		return NewHttpSpeechRecognizer(host, port, protocol)
	}

	return nil
}
