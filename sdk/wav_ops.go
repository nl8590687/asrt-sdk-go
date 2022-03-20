package sdk

import "github.com/nl8590687/asrt-sdk-go/common"

func LoadFile(filename string) []byte {
	return common.ReadBinFile(filename)
}

func DecodeWav(waveByte []byte) (*common.Wav, error) {
	wave := common.Wav{}
	err := wave.Deserialize(waveByte)
	if err != nil {
		return nil, err
	}

	return &wave, nil
}
