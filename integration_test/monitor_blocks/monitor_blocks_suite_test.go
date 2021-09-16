package monitor_blocks_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestMonitorBlocks(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "MonitorBlocks Suite")
}
