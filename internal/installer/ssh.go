package installer

import (
	"fmt"
	"github.com/mitchellh/go-homedir"
	"github.com/pkg/errors"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"sync"
	"time"
)

const (
	DefaultSshPort    = 22
	DefaultSshUser    = "root"
	DefaultSshKeyPath = "~/.ssh/id_rsa"
)

type Host struct {
	Address        string
	User           string
	Port           int8
	PrivateKeyPath string
}

func (h *Host) validate() error {
	if h.Address == "" {
		return errors.New("the address must be specified.")
	}
	if h.User == "" {
		h.User = DefaultSshUser
	}
	if h.Port <= 0 {
		h.Port = DefaultSshPort
	}
	if h.PrivateKeyPath == "" {
		h.PrivateKeyPath = DefaultSshKeyPath
	}
	return nil
}

func (h *Host) endpoint() string {
	return fmt.Sprintf("%s:%d", h.Address, h.Port)
}

type Connection struct {
	mux      sync.Mutex
	client  *ssh.Client
}

func (c *Connection) Close() {
	c.mux.Lock()
	defer c.mux.Unlock()

	if c.client == nil {
		return
	}

	if c.client != nil {
		c.client.Close()
		c.client = nil
	}
}

func (c *Connection) session() (*ssh.Session, error) {
	c.mux.Lock()
	defer c.mux.Unlock()

	if c.client == nil {
		return nil, errors.New("client is nil")
	}

	sess, err := c.client.NewSession()
	if err != nil {
		return nil, err
	}

	modes := ssh.TerminalModes{
		ssh.ECHO:          0,     // disable echoing
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}

	err = sess.RequestPty("xterm", 100, 50, modes)
	if err != nil {
		return nil, err
	}

	return sess, nil
}

func (c *Connection) Mkdir(dir string) error {
	session, err := c.session()
	if err != nil {
		return err
	}

	cmd := fmt.Sprintf("mkdir -p %s", dir)
	_, err = session.CombinedOutput(cmd)
	return err
}

func NewConnection(host *Host)(*Connection, error) {
	if err := host.validate(); err != nil {
		return nil, errors.Wrap(err, "validate Host failed")
	}

	config := &ssh.ClientConfig{
		Timeout:         time.Second * 5,
		User:            host.User,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	authMethod, err := publicKeyAuthMethod(host.PrivateKeyPath)
	if err != nil {
		return nil, err
	}
	config.Auth = []ssh.AuthMethod{authMethod}

	sshClient, err := ssh.Dial("tcp", host.endpoint(), config)
	if err != nil {
		return nil, errors.Wrap(err, "ssh dial failed")
	}

	return &Connection{
		client: sshClient,
	}, nil
}

func publicKeyAuthMethod(kPath string) (ssh.AuthMethod, error) {
	keyPath, err := homedir.Expand(kPath)
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadFile(keyPath)
	if err != nil {
		return nil, errors.Wrapf(err, "read public key file %s failed", kPath)
	}

	// Create the Signer for this private key.
	signer, err := ssh.ParsePrivateKey(data)
	if err != nil {
		return nil, errors.Wrap(err, "The given SSH key could not be parsed")
	}
	return ssh.PublicKeys(signer), nil
}
