package common_test

import (
	"fmt"
	"strings"

	"github.com/icon-project/btp/integration_test/setup"
	"github.com/onsi/gomega"
)

func Expect() {
	var res string
outer:
	for _, testcase := range BtpAddressTestCase {

		for i := range testcase.Output {

			for _, s := range i {
				if strings.Contains(s, "not supported blockchain") {
					re := strings.Split(s, ":")
					r := re[1] + ":" + re[2]

					res = strings.Trim(r, " ")
					fmt.Println(r)
					break outer

				}

			}
		}

	}
	gomega.Expect(res).To(gomega.Equal("not supported blockchain:xyz"))

}

var BtpAddressTestCase = []struct {
	Description string
	Input       struct {
		EnvironmentVariables []string
	}
	Output chan []string
}{
	{
		"Should fail on Invalid BTPAddress",
		struct{ EnvironmentVariables []string }{
			EnvironmentVariables: setup.NewEnvVariables(setup.EnvVariables{
				SrcAddress: "BTPSIMPLE_SRC_ADDRESS=btp://0x1.xyz/cx345676767788",
			}).ToValues(),
		},
		make(chan []string, 50),
	},
}
