package sieve

import (
	"bufio"
	"encoding/base64"
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"
)

type ManageSieve struct {
	conn    net.Conn
	scanner *bufio.Scanner

	addr string
}

type ManageSieveOpt func(*ManageSieve) error

func WithServerAddress(addr string) ManageSieveOpt {
	return func(c *ManageSieve) error {
		c.addr = addr
		return nil
	}
}

func WithConn(conn net.Conn) ManageSieveOpt {
	return func(c *ManageSieve) error {
		c.conn = conn
		return nil
	}
}

func NewManageSieve(opts ...ManageSieveOpt) (*ManageSieve, error) {
	c := &ManageSieve{addr: "localhost:4190"}

	for _, opt := range opts {
		if err := opt(c); err != nil {
			return nil, err
		}
	}

	conn, err := net.DialTimeout("tcp", c.addr, 5*time.Second)
	if err != nil {
		return nil, err
	}

	c.conn = conn
	c.scanner = bufio.NewScanner(c.conn)
	_, err = c.readResponse()

	return c, err
}

func (ms *ManageSieve) readResponse() ([]string, error) {
	var err error
	var res []string

	for ms.scanner.Scan() {
		line := strings.ToUpper(ms.scanner.Text())
		cmd := strings.ToUpper(line)
		if strings.HasPrefix(cmd, "OK") {
			err = ms.scanner.Err()
			break
		}
		if strings.HasPrefix(cmd, "NO") {
			err = errors.New(line[2:])
			break
		}
		if strings.HasPrefix(cmd, "BYE") {
			err = errors.New(line[3:])
			break
		}

		res = append(res, line)
	}
	return res, err
}

func (ms *ManageSieve) runCmd(cmd string, args ...string) ([]string, error) {
	for i, arg := range args {
		args[i] = strconv.Quote(arg)
	}
	_, _ = fmt.Fprint(ms.conn, cmd, " ", strings.Join(args, " "), "\r\n")

	return ms.readResponse()
}

// Login authenticates with managesieve server with given username and password,
// using PLAIN auth.
func (ms *ManageSieve) Login(user, pass string) error {
	auth := base64.StdEncoding.EncodeToString([]byte("\x00" + user + "\x00" + pass))
	_, err := ms.runCmd("AUTHENTICATE", "PLAIN", auth)
	return err
}

// GetScript gets sieve script by name.
func (ms *ManageSieve) GetScript(name string) (string, error) {
	s, err := ms.runCmd("GETSCRIPT", name)
	return strings.Join(s[2:], "\r\n"), err
}

// PutScript replace a sieve script with new content.
func (ms *ManageSieve) PutScript(name string, content string) error {
	content = fmt.Sprintf("{%d+}\r\n%s", len(content), content)
	_, err := ms.runCmd("PUTSCRIPT", name, content)
	return err
}

// SetActive marks the sieve script active.
func (ms *ManageSieve) SetActive(name string) error {
	_, err := ms.runCmd("SETACTIVE", name)
	return err
}

// DeleteScript deletes a sieve script by name.
func (ms *ManageSieve) DeleteScript(name string) error {
	_, err := ms.runCmd("DELETESCRIPT", name)
	return err
}
