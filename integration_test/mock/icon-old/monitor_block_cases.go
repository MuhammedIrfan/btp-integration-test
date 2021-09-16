package icon

import (
	"github.com/go-kit/kit/endpoint"
	"github.com/icon-project/btp/integration_test/setup/api"
	"github.com/onsi/gomega"
)

type UpdateMTATest1 struct {
	endpoints map[string]api.Handler
	output  chan string
}
type HelloParam struct {
	Name string `json:"name"`
}

func (i UpdateMTATest1) Expect() {
outer:
	for {
		select {
		case msg1 := <-i.output:
			gomega.Expect(msg1).To(gomega.Equal("irfan"))
			break outer

		}
	}
}

func (i UpdateMTATest1) Endpoints() map[string]api.Handler {
	return i.endpoints
}

func OutputCall(result interface{}, params interface{}, output chan interface{}) (func(ctx *api.Context, params *api.Params) (result interface{}, err error)) {
	return func(ctx *api.Context, p *api.Params) (result interface{}, err error) {
		if err := p.Convert(&params); err != nil {
			return nil, api.ErrInvalidParams()
		}
		output <- result
		close(output)
		return result, nil
	}
}

func GenericCall(result interface{}, params interface{}) (func(ctx *api.Context, params *api.Params) (result interface{}, err error)){
	return func(ctx *api.Context, p *api.Params) (result interface{}, err error) {
		if err := p.Convert(&params); err != nil {
			return nil, api.ErrInvalidParams()
		}
		return result, nil
	}
}


var UpdateMTATestCases = []struct {
	Description string
	Input       interface {
		Endpoints() map[string]api.Handler
		Expect()
	}
	output interface{}
}{
	{
		"should return 0x1",
		UpdateMTATest1{
			endpoints: map[string]api.Handler{
				"getStatus": GenericCall("test", ""),
				"hello": OutputCall("","", ),
			},
			output: make(chan string),
		},
		"0x1",
	},
}
