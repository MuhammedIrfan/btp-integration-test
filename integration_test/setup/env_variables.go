package setup

import "reflect"

type EnvVariables struct {
	BaseDir           string `default:"BTPSIMPLE_BASE_DIR=/btpsimple/data/btpsimple_src"`
	Config            string `default:"BTPSIMPLE_CONFIG=/btpsimple/config/src.config.json"`
	SrcAddress        string `default:"BTPSIMPLE_SRC_ADDRESS=btp://0x4e45e7.icon/cx419dddd0cb8d81d2e2229ec85c08f1ea1cefe1ff"`
	SrcEndpoint       string `default:"BTPSIMPLE_SRC_ENDPOINT=http://host.docker.internal:8081/api"`
	DstAddress        string `default:"BTPSIMPLE_DST_ADDRESS=btp://0x4d1014.icon/cx408d4c76cd46c55b19c1a2e978c88fd9e97faaac"`
	DstEndpoint       string `default:"BTPSIMPLE_DST_ENDPOINT=http://host.docker.internal:8080/api"`
	Offset            string `default:"BTPSIMPLE_OFFSET=3"`
	Keystore          string `default:"BTPSIMPLE_KEY_STORE=/btpsimple/config/src.ks.json"`
	KeySecret         string `default:"BTPSIMPLE_KEY_SECRET=/btpsimple/config/src.secret"`
	LogWriterFilename string `default:"BTPSIMPLE_LOG_WRITER_FILENAME=/btpsimple/data/log/btpsimple_src.log"`
}

func NewEnvVariables(e EnvVariables) EnvVariables {

	typ := reflect.TypeOf(e)

	if e.BaseDir == "" {

		f, _ := typ.FieldByName("BaseDir")

		e.BaseDir = f.Tag.Get("default")
	}
	if e.Config == "" {

		f, _ := typ.FieldByName("Config")

		e.Config = f.Tag.Get("default")
	}
	if e.SrcAddress == "" {

		f, _ := typ.FieldByName("SrcAddress")

		e.SrcAddress = f.Tag.Get("default")
	}
	if e.SrcEndpoint == "" {

		f, _ := typ.FieldByName("SrcEndpoint")

		e.SrcEndpoint = f.Tag.Get("default")
	}
	if e.DstAddress == "" {

		f, _ := typ.FieldByName("DstAddress")

		e.DstAddress = f.Tag.Get("default")
	}
	if e.DstEndpoint == "" {

		f, _ := typ.FieldByName("DstEndpoint")

		e.DstEndpoint = f.Tag.Get("default")
	}
	if e.Offset == "" {

		f, _ := typ.FieldByName("Offset")

		e.Offset = f.Tag.Get("default")
	}
	if e.Keystore == "" {

		f, _ := typ.FieldByName("Keystore")

		e.Keystore = f.Tag.Get("default")
	}
	if e.KeySecret == "" {

		f, _ := typ.FieldByName("KeySecret")

		e.KeySecret = f.Tag.Get("default")
	}
	if e.LogWriterFilename == "" {

		f, _ := typ.FieldByName("LogWriterFilename")

		e.LogWriterFilename = f.Tag.Get("default")
	}

	return e

}

func (e EnvVariables) ToValues() []string {

	v := reflect.ValueOf(e)
	res := []string{}

	for i := 0; i < v.NumField(); i++ {

		res = append(res, v.Field(i).String())
	}

	return res
}
