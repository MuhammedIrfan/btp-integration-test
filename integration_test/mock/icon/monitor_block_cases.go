package icon

import (
	"fmt"

	ty "github.com/icon-project/btp/cmd/btpsimple/module/icon"
	"github.com/icon-project/btp/integration_test/setup/api"
)

type UpdateMTATest struct {
	srcendpoints   map[string]api.Handler
	dstendpoints   map[string]api.Handler
	EndpointOutput struct {
		name   string
		params interface{}
		key    string
	}
	output chan string
}

type HelloParam struct {
	Name string `json:"name"`
}
type Genericall struct {
	Param interface{} `json:"params"`
}

type BMCdata struct {
	Method string      `json:"method"`
	Params interface{} `json:"params"`
}
type Functioncall struct {
	Height      string     `json:"height"`
	FromAddress ty.Address `json:"from" validate:"optional,t_addr_eoa"`
	ToAddress   ty.Address `json:"to" validate:"required,t_addr_score"`
	DataType    string     `json:"dataType" validate:"required,call"`
	Data        BMCdata    `json:"data"`
}

func (i UpdateMTATest) Expect() {
	// outer:
	for {
		select {
		case msg1 := <-output:

			// gomega.Expect(msg1).To(gomega.Equal("getStaus output"))
			// break outer
			fmt.Println(msg1)

		}
	}
}

func (i UpdateMTATest) SrcEndpoints() map[string]api.Handler {

	//i.endpoints[i.EndpointOutput.name] = i.OutputCall()

	return i.srcendpoints
}
func (i UpdateMTATest) DstEndpoints() map[string]api.Handler {

	//i.endpoints[i.EndpointOutput.name] = i.OutputCall()

	return i.dstendpoints
}

func (i UpdateMTATest) OutputCall() func(ctx *api.Context, params *api.Params) (interface{}, error) {
	return func(ctx *api.Context, p *api.Params) (interface{}, error) {
		var params Functioncall
		if err := p.Convert(&params); err != nil {
			return nil, api.ErrInvalidParams()
		}

		i.output <- ""
		close(i.output)

		return nil, nil
	}
}

func GenericCall(result interface{}, params interface{}) func(ctx *api.Context, params *api.Params) (interface{}, error) {

	return func(ctx *api.Context, p *api.Params) (interface{}, error) {
		fmt.Println("Genericcall")
		fmt.Println(params)
		if err := p.Convert(params); err != nil {
			return nil, api.ErrInvalidParams()
		}
		fmt.Println(result)
		output <- result
		return result, nil
	}
}

var output = make(chan interface{}, 50)

func FunctionCall(result interface{}, params interface{}, data interface{}, ap apis) func(ctx *api.Context, params *api.Params) (interface{}, error) {

	return func(ctx *api.Context, p *api.Params) (interface{}, error) {

		// var fn Functioncall
		if data == nil {
			res := params.(*Functioncall).Data.Method
			if err := p.Convert(params); err != nil {

				return nil, api.ErrInvalidParams()

			}
			r := ap[res]

			result = r

			output <- result
			fmt.Println("FunctionCall")

			return result, nil
		}

		if err := p.Convert(params); err != nil {

			return nil, api.ErrInvalidParams()

		}
		result = data
		output <- result

		return result, nil

	}
}

type BMCLinkStatus struct {
	TxSeq    ty.HexInt `json:"tx_seq"`
	RxSeq    ty.HexInt `json:"rx_seq"`
	Verifier struct {
		Height     ty.HexInt `json:"height"`
		Offset     ty.HexInt `json:"offset"`
		LastHeight ty.HexInt `json:"last_height"`
	} `json:"verifier"`
	BMRs []struct {
		Address      ty.Address `json:"address"`
		BlockCount   ty.HexInt  `json:"block_count"`
		MessageCount ty.HexInt  `json:"msg_count"`
	} `json:"relays"`
	BMRIndex         ty.HexInt `json:"relay_idx"`
	RotateHeight     ty.HexInt `json:"rotate_height"`
	RotateTerm       ty.HexInt `json:"rotate_term"`
	DelayLimit       ty.HexInt `json:"delay_limit"`
	MaxAggregation   ty.HexInt `json:"max_agg"`
	CurrentHeight    ty.HexInt `json:"cur_height"`
	RxHeight         ty.HexInt `json:"rx_height"`
	RxHeightSrc      ty.HexInt `json:"rx_height_src"`
	BlockIntervalSrc ty.HexInt `json:"block_interval_src"`
	BlockIntervalDst ty.HexInt `json:"block_interval_dst"`
}
type apis map[string]BMCLinkStatus

var param = &Functioncall{
	Data: BMCdata{
		Method: "getStatus",
	},
}

var param2 = &ty.BlockHeightParam{
	Height: ty.NewHexInt(0),
}
var param3 = &ty.DataHashParam{
	Hash: ty.NewHexBytes([]byte("1")),
}
var param4 = &ty.ProofResultParam{
	BlockHash: ty.NewHexBytes([]byte("1")),
	Index:     ty.NewHexInt(0),
}
var param5 = &ty.ProofEventsParam{
	BlockHash: ty.NewHexBytes([]byte("1")),
	Index:     ty.NewHexInt(0),
	Events: []ty.HexInt{
		ty.HexInt(1),
	},
}
var param6 = &ty.TransactionHashParam{
	Hash: ty.NewHexBytes([]byte("1")),
}
var param7 = &ty.TransactionParam{
	Version:     ty.NewHexInt(0),
	FromAddress: ty.NewAddress([]byte("1")),
	ToAddress:   ty.NewAddress([]byte("1")),
	Value:       ty.NewHexInt(0),
	StepLimit:   ty.NewHexInt(0),
	Timestamp:   ty.NewHexInt(0),
	NetworkID:   ty.NewHexInt(0),
	Nonce:       ty.NewHexInt(0),
	Signature:   "signature",
	DataType:    "datatype",
	Data:        "data",
	TxHash:      ty.NewHexBytes([]byte("1")),
}
var UpdateMTATestCases = []struct {
	Description string
	Input       interface {
		SrcEndpoints() map[string]api.Handler
		DstEndpoints() map[string]api.Handler
		Expect()
	}
	output interface{}
}{
	{
		"should return 0x1",
		UpdateMTATest{
			srcendpoints: map[string]api.Handler{
				"icx_getBlockHeaderByHeight": FunctionCall("", param2, "+G8CAAD4APgAoKhBUc6ulLLe8El4LGwCZ5Q0XjQLvxkLZKR+ae4JehAGoLyR7H03OhdjF+YHjNYGB6P7eHy6Fdv4E2ZTdcUpEhx8+ACgd46TPWgAPrZWZTwl9vSmLRvo9sJFR0eKg+JfJ6B4pf+A+AA=", apis{}),
				"icx_getDataByHash":          FunctionCall("", param3, "nil", apis{}),
				"icx_getVotesByHeight":       FunctionCall("", param2, "G8CAAD4APgAoKhBUc6ulLLe8El4LGwCZ5Q0XjQLvxkLZKR+ae4JehAGoLyR7H03OhdjF+YHjNYGB6P7eHy6Fdv4E2ZTdcUpEhx8+ACgd46TPWgAPrZWZTwl9vSmLRvo9sJFR0eKg+JfJ6B4pf+A+AA=", apis{}),
				"icx_getProofForResult":      FunctionCall("", param4, "nil", apis{}),
				"icx_getProofForEvents":      FunctionCall("", param5, "nil", apis{}),
				"icx_WaitTransactionResult":  FunctionCall("", param6, "nil", apis{}),
				"icx_GetTransactionResult":   FunctionCall("", param6, "nil", apis{}),
				"icx_SendTransactionAndWait": FunctionCall("", param7, "nil", apis{}),
				"icx_SendTransaction":        FunctionCall("", param7, "nil", apis{}),
				"icx_call": FunctionCall("", param, nil, apis{

					"getStatus": BMCLinkStatus{
						TxSeq: ty.NewHexInt(1),
						RxSeq: ty.NewHexInt(1),
						Verifier: struct {
							Height     ty.HexInt "json:\"height\""
							Offset     ty.HexInt "json:\"offset\""
							LastHeight ty.HexInt "json:\"last_height\""
						}{
							Height:     ty.NewHexInt(2871),
							Offset:     ty.NewHexInt(6),
							LastHeight: ty.NewHexInt(100),
						},
						BMRs: []struct {
							Address      ty.Address "json:\"address\""
							BlockCount   ty.HexInt  "json:\"block_count\""
							MessageCount ty.HexInt  "json:\"msg_count\""
						}{
							{
								Address:      "hx1133a33b41d65bcc02e148f6d61207194c97466b",
								BlockCount:   ty.NewHexInt(2865),
								MessageCount: ty.NewHexInt(1),
							},
						},
						BMRIndex:         ty.NewHexInt(0),
						RotateHeight:     ty.NewHexInt(2871),
						RotateTerm:       ty.NewHexInt(10),
						DelayLimit:       ty.NewHexInt(3),
						MaxAggregation:   ty.NewHexInt(10),
						CurrentHeight:    ty.NewHexInt(2885),
						RxHeight:         ty.NewHexInt(37),
						RxHeightSrc:      ty.NewHexInt(10),
						BlockIntervalSrc: ty.NewHexInt(1000),
						BlockIntervalDst: ty.NewHexInt(1000),
					},
				}),
			},
			dstendpoints: map[string]api.Handler{
				"icx_getBlockHeaderByHeight": FunctionCall("", param2, "+G8CAAD4APgAoKhBUc6ulLLe8El4LGwCZ5Q0XjQLvxkLZKR+ae4JehAGoLyR7H03OhdjF+YHjNYGB6P7eHy6Fdv4E2ZTdcUpEhx8+ACgd46TPWgAPrZWZTwl9vSmLRvo9sJFR0eKg+JfJ6B4pf+A+AA=", apis{}),
				"icx_getDataByHash":          FunctionCall("", param3, "nil", apis{}),
				"icx_getVotesByHeight":       FunctionCall("", param2, "G8CAAD4APgAoKhBUc6ulLLe8El4LGwCZ5Q0XjQLvxkLZKR+ae4JehAGoLyR7H03OhdjF+YHjNYGB6P7eHy6Fdv4E2ZTdcUpEhx8+ACgd46TPWgAPrZWZTwl9vSmLRvo9sJFR0eKg+JfJ6B4pf+A+AA=", apis{}),
				"icx_getProofForResult":      FunctionCall("", param4, "nil", apis{}),
				"icx_getProofForEvents":      FunctionCall("", param5, "nil", apis{}),
				"icx_WaitTransactionResult":  FunctionCall("", param6, "nil", apis{}),
				"icx_GetTransactionResult":   FunctionCall("", param6, "nil", apis{}),
				"icx_SendTransactionAndWait": FunctionCall("", param7, "nil", apis{}),
				"icx_SendTransaction":        FunctionCall("", param7, "nil", apis{}),
				"icx_call": FunctionCall("", param, nil, apis{

					"getStatus": BMCLinkStatus{
						TxSeq: ty.NewHexInt(1),
						RxSeq: ty.NewHexInt(1),
						Verifier: struct {
							Height     ty.HexInt "json:\"height\""
							Offset     ty.HexInt "json:\"offset\""
							LastHeight ty.HexInt "json:\"last_height\""
						}{
							Height:     ty.NewHexInt(2871),
							Offset:     ty.NewHexInt(6),
							LastHeight: ty.NewHexInt(1),
						},
						BMRs: []struct {
							Address      ty.Address "json:\"address\""
							BlockCount   ty.HexInt  "json:\"block_count\""
							MessageCount ty.HexInt  "json:\"msg_count\""
						}{
							{
								Address:      "hx1133a33b41d65bcc02e148f6d61207194c97466b",
								BlockCount:   ty.NewHexInt(2865),
								MessageCount: ty.NewHexInt(1),
							},
						},
						BMRIndex:         ty.NewHexInt(0),
						RotateHeight:     ty.NewHexInt(2871),
						RotateTerm:       ty.NewHexInt(10),
						DelayLimit:       ty.NewHexInt(3),
						MaxAggregation:   ty.NewHexInt(10),
						CurrentHeight:    ty.NewHexInt(2885),
						RxHeight:         ty.NewHexInt(37),
						RxHeightSrc:      ty.NewHexInt(10),
						BlockIntervalSrc: ty.NewHexInt(1000),
						BlockIntervalDst: ty.NewHexInt(1000),
					},
				}),
			},

			output: make(chan string),
			EndpointOutput: struct {
				name   string
				params interface{}
				key    string
			}{

				name:   "method",
				params: param,
				key:    "getStatus",
			},
		},
		"0x1",
	},
}
