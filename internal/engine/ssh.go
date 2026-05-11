package engine

import (
	"fmt"
	"io"
	"log"
	"net"
	"time"

	"golang.org/x/crypto/ssh"
)

// DelayedWriter wraps an io.Writer and adds delays between characters and lines
type DelayedWriter struct {
	w         io.Writer
	charDelay time.Duration
	lineDelay time.Duration
}

func (dw *DelayedWriter) Write(p []byte) (n int, err error) {
	for _, b := range p {
		nw, err := dw.w.Write([]byte{b})
		n += nw
		if err != nil {
			return n, err
		}
		if dw.charDelay > 0 {
			time.Sleep(dw.charDelay)
		}
		if dw.lineDelay > 0 && b == '\n' {
			time.Sleep(dw.lineDelay)
		}
	}
	return n, nil
}

// SSHServer represents an instance of an SSH server bound to a serial port
type SSHServer struct {
	Port       int
	SerialName string
	SerialConf SerialConfig
	Listener   net.Listener
	Config     *ssh.ServerConfig
	Quit       chan struct{}
}

// NewSSHServer creates a new SSH server configuration
func NewSSHServer(tcpPort int, serialName string, password string, serialConf SerialConfig) (*SSHServer, error) {
	config := &ssh.ServerConfig{
		PasswordCallback: func(c ssh.ConnMetadata, pass []byte) (*ssh.Permissions, error) {
			if string(pass) == password {
				return nil, nil
			}
			return nil, fmt.Errorf("password rejected for %q", c.User())
		},
	}

	return &SSHServer{
		Port:       tcpPort,
		SerialName: serialName,
		SerialConf: serialConf,
		Config:     config,
		Quit:       make(chan struct{}),
	}, nil
}

// Start starts the SSH server
func (s *SSHServer) Start(hostKey ssh.Signer) error {
	s.Config.AddHostKey(hostKey)

	addr := fmt.Sprintf(":%d", s.Port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	s.Listener = listener

	go func() {
		for {
			nConn, err := s.Listener.Accept()
			if err != nil {
				select {
				case <-s.Quit:
					return
				default:
					log.Printf("failed to accept incoming connection: %s", err)
					continue
				}
			}

			go s.handleConnection(nConn)
		}
	}()

	return nil
}

// Stop stops the SSH server
func (s *SSHServer) Stop() {
	close(s.Quit)
	if s.Listener != nil {
		s.Listener.Close()
	}
}

func (s *SSHServer) handleConnection(nConn net.Conn) {
	_, chans, reqs, err := ssh.NewServerConn(nConn, s.Config)
	if err != nil {
		log.Printf("failed to handshake: %s", err)
		return
	}

	// The incoming Request channel must be serviced.
	go ssh.DiscardRequests(reqs)

	for newChannel := range chans {
		if newChannel.ChannelType() != "session" {
			newChannel.Reject(ssh.UnknownChannelType, "unknown channel type")
			continue
		}
		channel, requests, err := newChannel.Accept()
		if err != nil {
			log.Printf("Could not accept channel: %v", err)
			continue
		}

		go func(in <-chan *ssh.Request) {
			for req := range in {
				switch req.Type {
				case "shell":
					req.Reply(true, nil)
				default:
					req.Reply(false, nil)
				}
			}
		}(requests)

		// Open serial port
		ser, err := OpenSerialPort(s.SerialName, s.SerialConf)
		if err != nil {
			log.Printf("failed to open serial port %s: %s", s.SerialName, err)
			channel.Close()
			continue
		}

		// Wrap serial port with DelayedWriter for transmit delay
		dw := &DelayedWriter{
			w:         ser,
			charDelay: time.Duration(s.SerialConf.CharDelay) * time.Millisecond,
			lineDelay: time.Duration(s.SerialConf.LineDelay) * time.Millisecond,
		}

		// Pipe data bi-directionally
		go func() {
			_, _ = io.Copy(channel, ser)
			channel.Close()
			ser.Close()
		}()
		go func() {
			_, _ = io.Copy(dw, channel)
			channel.Close()
			ser.Close()
		}()
	}
}
