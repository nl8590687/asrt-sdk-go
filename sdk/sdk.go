/*
 Copyright 2016-2099 Ailemon.net

 This file is part of Golang SDK ASRT Speech Recognition Tool.

 ASRT is free software: you can redistribute it and/or modify
 it under the terms of the GNU General Public License as published by
 the Free Software Foundation, either version 3 of the License, or
 (at your option) any later version.
 ASRT is distributed in the hope that it will be useful,
 but WITHOUT ANY WARRANTY; without even the implied warranty of
 MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 GNU General Public License for more details.

 You should have received a copy of the GNU General Public License
 along with ASRT.  If not, see <https://www.gnu.org/licenses/>.
 =====================================================================
*/

// package sdk ASRT语音识别接口调用SDK
package sdk

import "strings"

// GetSpeechRecognizer 获取一个语音识别调用类实例对象
func GetSpeechRecognizer(host string, port string, protocol string) ISpeechRecognizer {
	protocol = strings.ToLower(protocol)
	if protocol == "http" || protocol == "https" {
		return NewHTTPSpeechRecognizer(host, port, protocol, "")
	} else if protocol == "grpc" || protocol == "grpcs" {
		return NewGRPCSpeechRecognizer(host, port, protocol)
	}

	return nil
}
