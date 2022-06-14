package sdk

import (
	"context"
	"fmt"
	"log"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/nl8590687/asrt-sdk-go/common"
	grpcClient "github.com/nl8590687/asrt-sdk-go/grpc"
)

// NewGRPCSpeechRecognizer 调用ASRT语音识别系统HTTP+JSON协议接口的语音识别类
type GRPCSpeechRecognizer struct {
	BaseSpeechRecognizer
	Client     grpcClient.AsrtGrpcServiceClient
	connection *grpc.ClientConn
}

// NewGRPCSpeechRecognizer 构造一个用于调用grpc+pb协议接口的语音识别类实例对象
func NewGRPCSpeechRecognizer(host string, port string, protocol string) *GRPCSpeechRecognizer {
	protocol = strings.ToLower(protocol)
	if protocol != "grpc" && protocol != "grpcs" {
		return nil
	}

	base := BaseSpeechRecognizer{
		Host:     host,
		Port:     port,
		Protocol: protocol,
	}

	address := fmt.Sprintf("%s:%s", host, port)
	// 得到 gRPC 链接客户端句柄
	var conn *grpc.ClientConn
	var err error
	if protocol == "grpc" {
		conn, err = grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	} else {
		conn, err = grpc.Dial(address)
	}

	if err != nil {
		log.Printf("error: did not connect to `%s`, %s", address, err)
		return nil
	}

	grpcSpeechRecognizer := GRPCSpeechRecognizer{
		BaseSpeechRecognizer: base,
		// 将 proto 里面的服务句柄 和 gRPC句柄绑定
		Client:     grpcClient.NewAsrtGrpcServiceClient(conn),
		connection: conn,
	}

	return &grpcSpeechRecognizer
}

// Recognite 调用ASRT语音识别
func (g *GRPCSpeechRecognizer) Recognite(wavData []byte, frameRate int, channels int, byteWidth int,
) (*common.AsrtAPIResponse, error) {
	if len(wavData) > wavDataMaxLength {
		return nil, fmt.Errorf("error: %s `%d`, %s `%d`",
			"Too long wave sample byte length:", len(wavData),
			"the max length is", wavDataMaxLength)
	}
	grpcRequest := grpcClient.SpeechRequest{
		WavData: &grpcClient.WavData{
			Samples:    wavData,
			SampleRate: int32(frameRate),
			Channels:   int32(channels),
			ByteWidth:  int32(byteWidth),
		},
	}

	grpcResponse, err := g.Client.All(context.Background(), &grpcRequest, grpc.EmptyCallOption{})
	if err != nil {
		return nil, err
	}

	apiResponse := common.AsrtAPIResponse{
		StatusCode:    int(grpcResponse.StatusCode),
		StatucMesaage: grpcResponse.StatusMessage,
		Result:        grpcResponse.TextResult,
	}

	return &apiResponse, nil
}

// RecogniteSpeech 调用ASRT语音识别声学模型
func (g *GRPCSpeechRecognizer) RecogniteSpeech(wavData []byte, frameRate int, channels int, byteWidth int,
) (*common.AsrtAPIResponse, error) {
	if len(wavData) > wavDataMaxLength {
		return nil, fmt.Errorf("error: %s `%d`, %s `%d`",
			"Too long wave sample byte length:", len(wavData),
			"the max length is", wavDataMaxLength)
	}
	grpcRequest := grpcClient.SpeechRequest{
		WavData: &grpcClient.WavData{
			Samples:    wavData,
			SampleRate: int32(frameRate),
			Channels:   int32(channels),
			ByteWidth:  int32(byteWidth),
		},
	}

	grpcResponse, err := g.Client.Speech(context.Background(), &grpcRequest, grpc.EmptyCallOption{})
	if err != nil {
		return nil, err
	}

	apiResponse := common.AsrtAPIResponse{
		StatusCode:    int(grpcResponse.StatusCode),
		StatucMesaage: grpcResponse.StatusMessage,
		Result:        grpcResponse.ResultData,
	}

	return &apiResponse, nil
}

// RecogniteLanguage 调用ASRT语音识别语言模型
func (g *GRPCSpeechRecognizer) RecogniteLanguage(sequencePinyin []string) (*common.AsrtAPIResponse, error) {
	grpcRequest := grpcClient.LanguageRequest{
		Pinyins: sequencePinyin,
	}

	grpcResponse, err := g.Client.Language(context.Background(), &grpcRequest, grpc.EmptyCallOption{})
	if err != nil {
		return nil, err
	}

	apiResponse := common.AsrtAPIResponse{
		StatusCode:    int(grpcResponse.StatusCode),
		StatucMesaage: grpcResponse.StatusMessage,
		Result:        grpcResponse.TextResult,
	}

	return &apiResponse, nil
}

// RecogniteStream 调用ASRT语音识别来流式识别音频
func (g *GRPCSpeechRecognizer) RecogniteStream(wavChannel <-chan *common.Wav,
	resultChannel chan<- *common.AsrtAPIResponse,
) error {
	ctx := context.Background()
	streamClient, err := g.Client.Stream(ctx, grpc.EmptyCallOption{})
	if err != nil {
		return fmt.Errorf("error:%s", err.Error())
	}

	recvFunction := func(ctx context.Context) {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				grpcResponse, err := streamClient.Recv()
				if err != nil {
					log.Println("error: the function to recv stream asr result is error,", err.Error())
					close(resultChannel)
					return
				}

				apiResponse := &common.AsrtAPIResponse{
					StatusCode:    int(grpcResponse.StatusCode),
					StatucMesaage: grpcResponse.StatusMessage,
					Result:        grpcResponse.TextResult,
				}
				resultChannel <- apiResponse
			}
		}
	}
	go recvFunction(ctx)

	for value := range wavChannel {
		grpcRequest := grpcClient.SpeechRequest{
			WavData: &grpcClient.WavData{
				Samples:    value.GetRawSamples(),
				SampleRate: int32(value.FrameRate),
				Channels:   int32(value.Channels),
				ByteWidth:  int32(value.SampleWidth),
			},
		}
		err = streamClient.Send(&grpcRequest)
		if err != nil {
			return fmt.Errorf("error:%s", err.Error())
		}
	}

	err = streamClient.CloseSend()
	if err != nil {
		return fmt.Errorf("error:%s", err.Error())
	}

	_, cancel := context.WithCancel(ctx)
	cancel()

	return nil
}

// RecogniteLong 调用ASRT语音识别来识别长音频序列
func (g *GRPCSpeechRecognizer) RecogniteLong(wavData []byte, frameRate int, channels int, byteWidth int,
) ([]*common.AsrtAPIResponse, error) {
	if frameRate != 16000 {
		return nil, fmt.Errorf("error: unsupport wave sample rate `%d`", frameRate)
	}
	if channels != 1 {
		return nil, fmt.Errorf("error: unsupport wave channels number `%d`", channels)
	}
	if byteWidth != 2 {
		return nil, fmt.Errorf("error: unsupport wave byte width `%d`", byteWidth)
	}

	byteData := wavData
	var asrtResult []*common.AsrtAPIResponse
	duration := 2 * 16000 * 10

	index := 0
	for ; index < len(byteData)/duration+1; index++ {
		rsp, err := g.Recognite(byteData, frameRate, channels, byteWidth)
		if err != nil {
			return asrtResult, err
		}

		asrtResult = append(asrtResult, rsp)
	}

	return asrtResult, nil
}

// RecogniteFile 调用ASRT语音识别来识别指定文件名的音频文件
func (g *GRPCSpeechRecognizer) RecogniteFile(filename string) ([]*common.AsrtAPIResponse, error) {
	binData := common.ReadBinFile(filename)
	wavAudio := common.Wav{}
	err := wavAudio.Deserialize(binData)
	if err != nil {
		return nil, err
	}

	asrtResult, err := g.RecogniteLong(wavAudio.GetRawSamples(),
		wavAudio.FrameRate, wavAudio.Channels, wavAudio.SampleWidth)

	return asrtResult, err
}

func (g *GRPCSpeechRecognizer) Close() {
	g.connection.Close()
}
