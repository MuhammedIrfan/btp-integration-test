package config_test

import (
	"bufio"
	"context"
	"log"

	"github.com/icon-project/btp/integration_test/mock/icon"
	"github.com/icon-project/btp/integration_test/setup"
	. "github.com/onsi/ginkgo"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	ct "github.com/icon-project/btp/integration_test/test_cases/common_test"

	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
)

var _ = Describe("Parse Config", func() {
	var server setup.Server
	var resource *dockertest.Resource
	var err error
	cli, err := client.NewClientWithOpts()
	if err != nil {
		panic(err)
	}

	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}
	Describe("Loadin", func() {})

	Context("When starting the BMR ", func() {

		for _, testcase := range ct.BtpAddressTestCase {
			c := icon.Chain{}

			res := new(icon.Res)
			server = setup.NewServer(8080, c, res)

			server.Start()

			resource, err = pool.RunWithOptions(&dockertest.RunOptions{
				Repository: "btpsimple",
				Tag:        "latest",
				Env:        testcase.Input.EnvironmentVariables,
			}, func(hc *docker.HostConfig) {
				hc.ExtraHosts = append(hc.ExtraHosts, "host.docker.internal:host-gateway")

			})
			if err != nil {
				log.Fatalf("Could not start resource: %s", err)
			}

			go func() {
				reader, err := cli.ContainerLogs(context.Background(), resource.Container.ID, types.ContainerLogsOptions{
					ShowStdout: true,
					ShowStderr: true,
					Follow:     true,
					Timestamps: false,
				})
				if err != nil {
					panic(err)
				}

				scanner := bufio.NewScanner(reader)
				for scanner.Scan() {
					var res []string
					res = append(res, scanner.Text())

					testcase.Output <- res

					defer reader.Close()
				}
			}()

			It(testcase.Description, ct.Expect)

		}

	})

	AfterEach(func() {
		resource.Close()
		server.Stop()
	})

})
