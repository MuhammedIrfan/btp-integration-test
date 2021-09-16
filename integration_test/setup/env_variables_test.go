package setup

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefault(t *testing.T) {

	env := NewEnvVariables(EnvVariables{})

	assert.Equal(t, "BTPSIMPLE_OFFSET=0", env.Offset)

}
func TestOffset(t *testing.T) {

	env := NewEnvVariables(EnvVariables{
		Offset: "BTPSIMPLE_OFFSET=5",
	})

	assert.Equal(t, "BTPSIMPLE_OFFSET=5", env.Offset)

}

func TestToArray(t *testing.T) {

	env := NewEnvVariables(EnvVariables{
		Offset: "BTPSIMPLE_OFFSET=5",
	})

	res := env.ToValues()

	assert.Equal(t, []string{
		"BTPSIMPLE_BASE_DIR=/btpsimple/data/btpsimple_src",
		"BTPSIMPLE_CONFIG=/btpsimple/config/src.config.json",
		"BTPSIMPLE_SRC_ADDRESS=btp://0x1.icon/cx345676767788",
		"BTPSIMPLE_SRC_ENDPOINT=http://host.docker.internal:8080/api",
		"BTPSIMPLE_DST_ADDRESS=btp://0x1.icon/cx345676767788",
		"BTPSIMPLE_DST_ENDPOINT=http://host.docker.internal:8080/api",
		"BTPSIMPLE_OFFSET=5",
		"BTPSIMPLE_KEY_STORE=/btpsimple/config/src.ks.json",
		"BTPSIMPLE_KEY_SECRET=/btpsimple/config/src.secret",
		"BTPSIMPLE_LOG_WRITER_FILENAME=/btpsimple/data/log/btpsimple_src.log",
	}, res)

}
