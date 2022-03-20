package installer

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"time"
)

const (
	SshPort        = 22
	SshDefaultUser = "root"
)

func NewSshSession(host string) (*ssh.Session, error) {
	config := &ssh.ClientConfig{
		Timeout:         time.Second * 5,
		User:            SshDefaultUser,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	addr := fmt.Sprintf("%s:%d", host, SshPort)
	sshClient, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		return nil, err
	}
	defer sshClient.Close()

	session, err := sshClient.NewSession()
	if err != nil {
		return nil, err
	}
	defer session.Close()
	return session, err
}

func CreateRemoteDir(host, dir string) error {
	session, err := NewSshSession(host)
	if err != nil {
		return err
	}
	cmd := fmt.Sprintf("mkdir -p %s", dir)
	_, err = session.CombinedOutput(cmd)
	return err
}
