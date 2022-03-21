package sdk

import "github.com/nl8590687/asrt-sdk-go/common"

// LoadFile 加载二进制文件
func LoadFile(filename string) []byte {
	return common.ReadBinFile(filename)
}

// DecodeWav 从wave格式byte数组反序列化解码为Wav对象
func DecodeWav(waveByte []byte) (*common.Wav, error) {
	wave := common.Wav{}
	err := wave.Deserialize(waveByte)
	if err != nil {
		return nil, err
	}

	return &wave, nil
}
