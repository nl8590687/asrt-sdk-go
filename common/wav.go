package common

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

// Wav Wav格式结构对象
type Wav struct {
	// Samples 采样样本数据
	Samples [][]int16
	// FrameRate 采样频率，单位：Hz。例如：16000 / 8000 等
	FrameRate int
	// Channels 声音通道数，单声道为1，立体声为2
	Channels int
	// SampleWidth 采样位深，单位：字节(byte)
	SampleWidth int
	// BytesPerSec 比特率，单位：bps
	BytesPerSec int
	wavByteData []byte

	blocklenSample uint16
	bitNum         uint16
	// fmtHeadLength  uint32
	// fmtHeader      []byte

	cksize uint32
}

// NewBlankWav 获取一个新的空白Wav对象
func NewBlankWav(frameRate int, channels int, sampleWidth int) Wav {
	wave := Wav{
		FrameRate:   frameRate,
		Channels:    channels,
		SampleWidth: sampleWidth,
		BytesPerSec: frameRate * channels * sampleWidth,
		Samples:     make([][]int16, channels),
	}

	for i := 0; i < channels; i += 1 {
		wave.Samples[i] = make([]int16, 0)
	}
	return wave
}

// Deserialize Wave格式反序列化
func (w *Wav) Deserialize(bytesData []byte) error {
	if bytesData == nil {
		return fmt.Errorf("error: byte array is nil")
	}

	w.wavByteData = bytesData
	bodyLength, p, err := w.parseHeader()
	if err != nil {
		return err
	}

	err = w.parseBody(p, bodyLength)
	return err
}

// parseHeader 解码头部
func (w *Wav) parseHeader() (bodyLength uint32, startPosition uint32, err error) {
	var riff uint32       // 4 byte
	var riffSize uint32   // 4 byte
	var waveID uint32     // 4 byte
	var junklength uint32 // 4 byte

	// var fmtID uint32  // 4 byte

	var cksize uint32 // 4 byte
	var waveType int
	var channel uint16     // 2 byte
	var sampleRate uint32  // 4 byte
	var bytespersec uint32 // 4 byte
	// var blocklenSample uint16 // 2 byte
	// var bitNum uint16         // 2 byte
	// var unknown uint16        // 2 byte
	// var dataID [4]byte
	// var dataLength uint32 // 4 byte

	var p uint32 = 0

	riff = binary.BigEndian.Uint32((w.wavByteData)[p : p+4])
	p += 4
	if riff != 0x52494646 {
		return 0, p, fmt.Errorf("error: this file is not riff format")
	}

	riffSize = binary.LittleEndian.Uint32((w.wavByteData)[p : p+4]) // 文件剩余长度
	p += 4
	if riffSize != uint32(len(w.wavByteData))-p {
		return 0, p, fmt.Errorf("error: this file maybe has been destroyed so that file length not equals flag value")
	}

	waveID = binary.BigEndian.Uint32((w.wavByteData)[p : p+4]) // wave文件标识
	p += 4
	if waveID != 0x57415645 {
		return 0, p, fmt.Errorf("error: this file is not wave file")
	}

	tmp := binary.BigEndian.Uint32((w.wavByteData)[p : p+4]) // 4 byte
	p += 4
	switch tmp {
	case 0x4A554E4B: // 发现了junk flag，这个值是 junkID
		junklength = binary.LittleEndian.Uint32((w.wavByteData)[p : p+4]) // junk长度
		p += 4
		p += junklength // 将不要的junk部分跳过

		_ = binary.BigEndian.Uint32((w.wavByteData)[p : p+4]) // 读fmt 标记: fmtID
		p += 4
	case 0x666D7420: // 发现了fmt flag，这个值是 fmtID
		_ = tmp // fmtID
	default:
		return 0, p, fmt.Errorf("error: can not find any junk or fmt flag in this wave file")
	}

	w.cksize = binary.LittleEndian.Uint32((w.wavByteData)[p : p+4]) // 4 byte，小端存储
	p += 4
	pDataStart := cksize
	_ = pDataStart + 8

	tmpWaveType := binary.LittleEndian.Uint16((w.wavByteData)[p : p+2]) // 2 byte，这个字段是小端存储
	p += 2
	waveType = int(tmpWaveType)
	if waveType != 1 {
		return 0, p, fmt.Errorf("error: this wave file is not pcm format and it is not supported")
	}

	channel = binary.LittleEndian.Uint16((w.wavByteData)[p : p+2]) // 声道数 2 byte，小端存储
	p += 2
	w.Channels = int(channel)

	sampleRate = binary.LittleEndian.Uint32((w.wavByteData)[p : p+4]) // 采样频率，小端存储
	p += 4
	w.FrameRate = int(sampleRate)

	bytespersec = binary.LittleEndian.Uint32((w.wavByteData)[p : p+4]) // 每秒钟字节数，小端存储
	p += 4
	w.BytesPerSec = int(bytespersec)

	w.blocklenSample = binary.LittleEndian.Uint16((w.wavByteData)[p : p+2]) // 每次采样的字节大小，2为单声道，4为立体声道，小端存储
	p += 2

	w.bitNum = binary.LittleEndian.Uint16((w.wavByteData)[p : p+2]) // 每个声道的采样精度，默认16bit，小端存储
	w.SampleWidth = int(w.bitNum) / 8
	p += 2

	tmp1 := binary.BigEndian.Uint16((w.wavByteData)[p : p+2])
	p += 2
	for tmp1 != 0x6461 { // 寻找da标记
		tmp1 = binary.BigEndian.Uint16((w.wavByteData)[p : p+2])
		p += 2
	}
	tmp1 = binary.BigEndian.Uint16((w.wavByteData)[p : p+2])
	p += 2
	if tmp1 != 0x7461 { // ta标记
		return 0, p, fmt.Errorf("error: can not find `data` flag in wave file")
	}

	dataSize := binary.LittleEndian.Uint32((w.wavByteData)[p : p+4]) // wav数据byte长度，小端存储
	p += 4
	if dataSize < 2 {
		dataSize = uint32(len(w.wavByteData)) - p
	}

	return dataSize, p, nil
}

// parseBody 解码数据区
func (w *Wav) parseBody(startPosition uint32, bodyLength uint32) error {
	p := startPosition
	numSamples := bodyLength / uint32(w.blocklenSample) // 计算样本数

	w.Samples = make([][]int16, w.Channels, numSamples)
	for i := 0; i < int(numSamples); i += 1 {
		for j := 0; j < w.Channels; j += 1 {
			// 读入2字节有符号整数
			bytesBuffer := bytes.NewBuffer(w.wavByteData[p : p+2])
			var tmpSample int16
			binary.Read(bytesBuffer, binary.LittleEndian, &tmpSample)
			p += uint32(w.bitNum / 8)
			w.Samples[j] = append(w.Samples[j], tmpSample)
		}
	}

	return nil
}

// Serialize Wave格式序列化
func (w *Wav) Serialize() ([]byte, error) {
	waveData := w.packWave()

	res, err := w.packRiff(waveData)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// packWave 打包wave部分
func (w *Wav) packWave() []byte {
	var byteBuf bytes.Buffer
	var tmpBytes []byte
	// waveID
	tmpBytes = make([]byte, 4)
	binary.BigEndian.PutUint32(tmpBytes, 0x57415645)
	byteBuf.Write(tmpBytes)
	// fmt flag
	binary.BigEndian.PutUint32(tmpBytes, 0x666D7420)
	byteBuf.Write(tmpBytes)
	// fmt header length
	binary.LittleEndian.PutUint32(tmpBytes, 0x00000010)
	byteBuf.Write(tmpBytes)
	// wave type: pcm
	tmpBytes = make([]byte, 2)
	binary.LittleEndian.PutUint16(tmpBytes, 0x0001)
	byteBuf.Write(tmpBytes)
	// channel count
	binary.LittleEndian.PutUint16(tmpBytes, uint16(w.Channels))
	byteBuf.Write(tmpBytes)
	// sample rate
	tmpBytes = make([]byte, 4)
	binary.LittleEndian.PutUint32(tmpBytes, uint32(w.FrameRate))
	byteBuf.Write(tmpBytes)
	// bytes per second
	binary.LittleEndian.PutUint32(tmpBytes, uint32(w.BytesPerSec))
	byteBuf.Write(tmpBytes)
	// block length sample
	tmpBytes = make([]byte, 2)
	binary.LittleEndian.PutUint16(tmpBytes, uint16(w.Channels*w.SampleWidth))
	byteBuf.Write(tmpBytes)
	// sample width
	binary.LittleEndian.PutUint16(tmpBytes, uint16(w.SampleWidth*8))
	byteBuf.Write(tmpBytes)
	// flag: data
	tmpBytes = make([]byte, 4)
	binary.BigEndian.PutUint32(tmpBytes, 0x64617461)
	byteBuf.Write(tmpBytes)
	// data length (byte)
	binary.LittleEndian.PutUint32(tmpBytes, uint32(len(w.Samples)*len(w.Samples[0])*2))
	byteBuf.Write(tmpBytes)

	// wave data
	tmpBytes = make([]byte, 2)
	for j := 0; j < len(w.Samples[0]); j += 1 {
		for i := 0; i < len(w.Samples); i += 1 {
			binary.LittleEndian.PutUint16(tmpBytes, uint16(w.Samples[i][j]))
			byteBuf.Write(tmpBytes)
		}
	}

	return byteBuf.Bytes()
}

// GetRawSamples 读取Wave格式的Samples原始数据
func (w *Wav) GetRawSamples() []byte {
	var byteBuf bytes.Buffer
	// wave data
	tmpBytes := make([]byte, 2)
	for j := 0; j < len(w.Samples[0]); j += 1 {
		for i := 0; i < len(w.Samples); i += 1 {
			binary.LittleEndian.PutUint16(tmpBytes, uint16(w.Samples[i][j]))
			byteBuf.Write(tmpBytes)
		}
	}
	return byteBuf.Bytes()
}

// packRiff 打包wave数据为riff格式文件
func (w *Wav) packRiff(waveData []byte) ([]byte, error) {
	if len(waveData) == 0 {
		return nil, fmt.Errorf("error: wave data sequence is empty")
	}

	var byteBuf bytes.Buffer
	var tmpBytes []byte

	// RIFF
	tmpBytes = make([]byte, 4)
	binary.BigEndian.PutUint32(tmpBytes, 0x52494646)
	byteBuf.Write(tmpBytes)
	// Length
	tmpBytes = make([]byte, 4)
	binary.LittleEndian.PutUint32(tmpBytes, uint32(len(waveData)))
	byteBuf.Write(tmpBytes)
	// Wave sequence
	byteBuf.Write(waveData)
	return byteBuf.Bytes(), nil
}

// AppendWav 在本wave的samples后面追加给定wave的sample
func (w *Wav) AppendWav(wavAppended Wav) error {
	if w.FrameRate != wavAppended.FrameRate {
		return fmt.Errorf(
			"error: appended wav's frame rate not equals this wav's. this wav's is %d but apeended wav's is %d",
			w.FrameRate, wavAppended.FrameRate)
	}
	if w.Channels != wavAppended.Channels {
		return fmt.Errorf(
			"error: appended wav's channel count not equals this wav's. this wav's is %d but apeended wav's is %d",
			w.Channels, wavAppended.Channels)
	}
	if w.SampleWidth != wavAppended.SampleWidth {
		return fmt.Errorf(
			"error: appended wav's sample width not equals this wav's. this wav's is %d but apeended wav's is %d",
			w.SampleWidth, wavAppended.SampleWidth)
	}
	if len(w.Samples) != len(wavAppended.Samples) {
		return fmt.Errorf("error: appended wav samples's shape not equals this samples's")
	}
	if len(w.Samples) == 0 || len(wavAppended.Samples) == 0 {
		return fmt.Errorf("error: wav samples's shape is zero")
	}

	for i := 0; i < len(w.Samples); i += 1 {
		w.Samples[i] = append(w.Samples[i], wavAppended.Samples[i]...)
	}

	return nil
}

/* AppendBlank 在wave的后面追加一定时间的静音区，单位：毫秒。
这里用无符号类型是因为不允许添加负数时间长度 */
func (w *Wav) AppendBlank(millisecond uint32) {
	zeroByteCount := millisecond * uint32(w.FrameRate/1000)
	for i := 0; i < w.Channels; i += 1 {
		arr1 := make([]int16, zeroByteCount)
		w.Samples[i] = append(w.Samples[i], arr1...)
	}
}
