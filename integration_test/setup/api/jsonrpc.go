package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"reflect"
	"strings"

	"github.com/labstack/echo/v4"
)

const Version = "2.0"

type Request struct {
	Version string          `json:"jsonrpc" validate:"required,version"`
	Method  string          `json:"method" validate:"required"`
	Params  json.RawMessage `json:"params,omitempty"`
	ID      interface{}     `json:"id"`
}

type Response struct {
	Version string      `json:"jsonrpc"`
	Result  interface{} `json:"result"`
	Error   *Error      `json:"error,omitempty"`
	ID      interface{} `json:"id"`
}

type ErrorResponse struct {
	Version string      `json:"jsonrpc"`
	Error   *Error      `json:"error"`
	ID      interface{} `json:"id"`
}

type Context struct {
	echo.Context
}

func NewContext(c echo.Context) *Context {
	ctx := &Context{Context: c}
	return ctx
}

type Params struct {
	rawMessage json.RawMessage
	validator  echo.Validator
}

func (p *Params) Convert(v interface{}) error {
	if p.rawMessage == nil {
		return errors.New("params message is null")
	}
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return errors.New("v is not pointer type or v is nil")
	}
	jd := json.NewDecoder(bytes.NewBuffer(p.rawMessage))
	jd.DisallowUnknownFields()
	if err := jd.Decode(v); err != nil {
		return err
	}
	return nil
}

func (p *Params) RawMessage() []byte {
	bs, _ := p.rawMessage.MarshalJSON()
	return bs
}

func (p *Params) IsEmpty() bool {

	return p.rawMessage == nil
}

func JsonRpc(mr *MethodRepository) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ctype := c.Request().Header.Get(echo.HeaderContentType)
			if !strings.HasPrefix(ctype, echo.MIMEApplicationJSON) {
				c.Request().Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			}
			r := new(Request)
			if err := c.Bind(r); err != nil {
				return ErrParse()
			}
			c.Set("request", r)
			method, err := mr.TakeMethod(r)
			if err != nil {
				return err
			}
			c.Set("method", method)
			return next(c)
		}
	}
}
