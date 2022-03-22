package test

import (
	"github.com/DataWorkbench/deploy/internal/installer"
	"testing"
)

const TestHostAddress = "worker-p004"

func TestSsh(t *testing.T) {
	host := &installer.Host{
		Address: TestHostAddress,
	}

	sshConn, err := installer.NewConnection(host)
	if err != nil {
		t.Fatalf("new ssh connection error: %+v", err)
	}
	defer sshConn.Close()

	err = sshConn.Mkdir("/tmp/test_dir")
	if err != nil {
		t.Fatalf("mkdir error: %+v", err)
	}
}
