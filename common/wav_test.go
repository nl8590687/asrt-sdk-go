package common

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"testing"

	"github.com/stretchr/testify/suite"
)

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestUnitWav(t *testing.T) {
	suite.Run(t, new(TestUnitWavSuite))
}

// Define the suite, and absorb the built-in basic suite
// functionality from testify - including a T() method which
// returns the current testing context
type TestUnitWavSuite struct {
	suite.Suite
	VariableThatShouldStartAtFive int
}

// Make sure that VariableThatShouldStartAtFive is set to five
// before each test
func (t *TestUnitWavSuite) SetupTest() {
	t.VariableThatShouldStartAtFive = 5
}

func (t *TestUnitWavSuite) TestNewBlankWav() {
	wave := NewBlankWav(16000, 1, 2)
	t.Equal(16000, wave.FrameRate)
	t.Equal(1, wave.Channels)
	t.Equal(2, wave.SampleWidth)
	t.Equal(32000, wave.BytesPerSec)

	wb := readBinFile("../testData/data1.wav")
	w := Wav{}
	err := w.Deserialize(wb)
	t.Equal(true, err == nil)

	wave.AppendWav(w)
	wb2, err := wave.Serialize()
	t.Equal(true, err == nil)

	err = writeBinFile("../testData/tmp.wav", wb2)
	t.Equal(true, err == nil)

	w2 := Wav{}
	err = w2.Deserialize(wb2)
	t.Equal(true, err == nil)
}

func (t *TestUnitWavSuite) TestDeserialize() {
	tests := []struct {
		name     string
		waveFile string
		want     bool
	}{
		{
			name:     "success",
			waveFile: "../testData/data1.wav",
			want:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func() {
			wavBytes := readBinFile(tt.waveFile)
			fmt.Println("waveBytes的长度:", len(wavBytes))

			wave := Wav{}
			err := wave.Deserialize(wavBytes)
			fmt.Println("wav类:", wave.BytesPerSec, wave.Channels, wave.FrameRate, wave.SampleWidth, wave.Samples[0][0:100])

			log.Println(err)
			t.Equal(tt.want, err == nil)
		})
	}
}

func (t *TestUnitWavSuite) TestSerialize() {
	tests := []struct {
		name     string
		waveFile string
		want     bool
	}{
		{
			name:     "success",
			waveFile: "../testData/data1.wav",
			want:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func() {
			wavBytes := readBinFile(tt.waveFile)
			fmt.Println("waveBytes的长度:", len(wavBytes))

			wave := Wav{}
			_ = wave.Deserialize(wavBytes)
			fmt.Println("wav类:", wave.BytesPerSec, wave.Channels, wave.FrameRate, wave.SampleWidth, wave.Samples[0][0:100])

			waveBytesNew, _ := wave.Serialize()
			_ = writeBinFile("../testData/tmp.wav", waveBytesNew)

			wavBytesNew2 := readBinFile(tt.waveFile)
			fmt.Println("waveBytes的长度:", len(wavBytesNew2))

			waveNew := Wav{}
			err := waveNew.Deserialize(wavBytesNew2)

			log.Println(err)
			t.Equal(tt.want, err == nil)
		})
	}
}

func (t *TestUnitWavSuite) TestAppendWav() {
	tests := []struct {
		name     string
		waveFile string
		want     bool
	}{
		{
			name:     "success",
			waveFile: "../testData/data1.wav",
			want:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func() {
			wavBytes := readBinFile(tt.waveFile)
			fmt.Println("waveBytes的长度:", len(wavBytes))

			wave1 := Wav{}
			err := wave1.Deserialize(wavBytes)
			t.Equal(tt.want, err == nil)
			fmt.Println("wav类:", wave1.BytesPerSec, wave1.Channels, wave1.FrameRate, wave1.SampleWidth, wave1.Samples[0][0:100])

			wave2 := Wav{}
			err = wave2.Deserialize(wavBytes)
			t.Equal(tt.want, err == nil)

			err = wave1.AppendWav(wave2)
			t.Equal(tt.want, err == nil)

			waveBytesNew, err := wave1.Serialize()
			t.Equal(tt.want, err == nil)
			err = writeBinFile("../testData/tmp.wav", waveBytesNew)
			t.Equal(tt.want, err == nil)
		})
	}
}

func (t *TestUnitWavSuite) TestAppendBlank() {
	tests := []struct {
		name     string
		waveFile string
		want     bool
	}{
		{
			name:     "success",
			waveFile: "../testData/data1.wav",
			want:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func() {
			wavBytes := readBinFile(tt.waveFile)
			fmt.Println("waveBytes的长度:", len(wavBytes))

			wave1 := Wav{}
			err := wave1.Deserialize(wavBytes)
			t.Equal(tt.want, err == nil)
			fmt.Println("wav类:", wave1.BytesPerSec, wave1.Channels, wave1.FrameRate, wave1.SampleWidth, wave1.Samples[0][0:100])

			wave2 := Wav{}
			err = wave2.Deserialize(wavBytes)
			t.Equal(tt.want, err == nil)

			wave1.AppendBlank(1000)
			err = wave1.AppendWav(wave2)
			t.Equal(tt.want, err == nil)

			waveBytesNew, err := wave1.Serialize()
			t.Equal(tt.want, err == nil)
			err = writeBinFile("../testData/tmp.wav", waveBytesNew)
			t.Equal(tt.want, err == nil)
		})
	}
}

func TestUnitDefault(t *testing.T) {
	var byteBuf bytes.Buffer
	var tmpBytes []byte
	tmpBytes = make([]byte, 4)
	binary.BigEndian.PutUint32(tmpBytes, 0x52494646)
	fmt.Println("11111111:", tmpBytes)
	byteBuf.Write(tmpBytes)
	fmt.Println("byteBuf", byteBuf.Bytes())
	tmpBytes = make([]byte, 2)
	fmt.Println("22222222:", tmpBytes)
}
