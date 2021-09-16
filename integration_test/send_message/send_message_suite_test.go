package send_message_test

import (
	"fmt"
	"os/exec"

	c "github.com/icon-project/btp/integration_test/send_message"
	s "github.com/icon-project/btp/integration_test/setup"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestSendMessage(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "SendMessage Suite")
}

var _ = Describe("test", func() {

	BeforeEach(func() {
		cmd := exec.Command("go", "run", "main.go", "&")
		cmd.Dir = "../"
		out, err := cmd.Output()
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(out)
	})
	Context("testing ", func() {
		It("actual test", func() {

			res := new(c.Res)
			ch := c.Chhain{}

			r := s.SetupRoutes(ch, res)
			fmt.Println(r)
			Expect(r).To(Equal("hi"))

		})
	})
})
