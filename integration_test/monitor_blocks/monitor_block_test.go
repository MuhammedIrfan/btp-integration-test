package monitor_blocks_test

import (
	"log"

	// "github.com/gorilla/websocket"

	"github.com/icon-project/btp/integration_test/mock/icon"
	"github.com/icon-project/btp/integration_test/setup"
	. "github.com/onsi/ginkgo"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	//. "github.com/onsi/gomega"
)

var _ = Describe("Block Monitoring", func() {
	var serversrc setup.Server
	var serverdst setup.Server

	var _ *dockertest.Resource
	var err error
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}
	Describe("loading from JSON", func() {})
	BeforeEach(func() {
	})

	Context("when BMR start montior", func() {
		for _, testcase := range icon.UpdateMTATestCases {

			c := icon.Chain{}

			res := new(icon.Res)
			serversrc = setup.NewServer(8080, c, res)
			serverdst = setup.NewServer(8081, c, res)

			// for k, v := range testcase.Input.Endpoints() {
			// 	server.RegisterMethods(map[string]api.Handler{
			// 		k : v,
			// 	})
			// }
			serversrc.RegisterMethods(testcase.Input.SrcEndpoints())
			serverdst.RegisterMethods(testcase.Input.DstEndpoints())
			serversrc.Start()
			serverdst.Start()

			_, err = pool.RunWithOptions(&dockertest.RunOptions{
				Repository:   "btpsimple",
				Tag:          "latest",
				ExposedPorts: []string{"40000"},
				PortBindings: map[docker.Port][]docker.PortBinding{
					"40000": {
						{HostIP: "", HostPort: "10080"},
					},
				},
				ExtraHosts: []string{"host.docker.internal:host-gateway"},
				Env:        setup.NewEnvVariables(setup.EnvVariables{}).ToValues(),
			})
			if err != nil {
				log.Fatalf("Could not start resource: %s", err)
			}
			// time.Sleep(100 * time.Second)

			It(testcase.Description, testcase.Input.Expect)
		}

		// It("should update the internal DB MTA to heightMTA at 100", func() {
		// 	ws, _, err := websocket.DefaultDialer.Dial("ws://localhost:8080/block", nil)
		// 	if err != nil {
		// 		log.Fatalf("%v", err)
		// 	}
		// 	defer ws.Close()

		// 	if err := ws.WriteMessage(websocket.TextMessage, []byte("hi")); err != nil {
		// 		log.Fatalf("%v", err)
		// 	}
		// 	_, p, err := ws.ReadMessage()
		// 	if err != nil {
		// 		log.Fatalf("%v", err)
		// 	}
		// 	time.Sleep(100 * time.Second)
		// 	Expect(string(p)).To(Equal("hello"))
		// })

	})

	AfterEach(func() {
		//resource.Close()
		//server.Stop()
	})
})
