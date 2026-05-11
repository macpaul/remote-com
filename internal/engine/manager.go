package engine

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"

	"golang.org/x/crypto/ssh"
)

// Binding represents a mapping between a serial port and a TCP port
type Binding struct {
	SerialPort string       `json:"serialPort"`
	TCPPort    int          `json:"tcpPort"`
	Password   string       `json:"password"`
	SerialConf SerialConfig `json:"serialConf"`
	Active     bool         `json:"active"`
}

// Manager handles the lifecycle of serial-ssh bindings
type Manager struct {
	Bindings map[string]*Binding // Key is SerialPort:TCPPort
	Servers  map[string]*SSHServer
	mu       sync.RWMutex
	hostKey  ssh.Signer
	ConfigPath string
}

func NewManager(configPath string) (*Manager, error) {
	m := &Manager{
		Bindings:   make(map[string]*Binding),
		Servers:    make(map[string]*SSHServer),
		ConfigPath: configPath,
	}

	// Load or generate host key
	key, err := m.getOrGenerateHostKey()
	if err != nil {
		return nil, err
	}
	m.hostKey = key

	// Load config if exists
	_ = m.LoadConfig()

	return m, nil
}

func (m *Manager) getOrGenerateHostKey() (ssh.Signer, error) {
	keyPath := "host_key.pem"
	keyData, err := os.ReadFile(keyPath)
	if err != nil {
		// Generate new key
		privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
		if err != nil {
			return nil, err
		}
		
		privateKeyPEM := &pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
		}
		
		f, err := os.Create(keyPath)
		if err != nil {
			return nil, err
		}
		defer f.Close()
		
		if err := pem.Encode(f, privateKeyPEM); err != nil {
			return nil, err
		}
		
		return ssh.NewSignerFromKey(privateKey)
	}

	return ssh.ParsePrivateKey(keyData)
}

func (m *Manager) LoadConfig() error {
	data, err := os.ReadFile(m.ConfigPath)
	if err != nil {
		return err
	}
	
	lines := strings.Split(string(data), "\n")
	var currentKey string
	var currentBinding *Binding
	
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Bindings = make(map[string]*Binding)
	
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, ";") || strings.HasPrefix(line, "#") {
			continue
		}
		
		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			currentKey = line[1 : len(line)-1]
			currentBinding = &Binding{}
			m.Bindings[currentKey] = currentBinding
			continue
		}
		
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 || currentBinding == nil {
			continue
		}
		
		key := strings.TrimSpace(parts[0])
		val := strings.TrimSpace(parts[1])
		
		switch key {
		case "SerialPort":
			currentBinding.SerialPort = val
		case "TCPPort":
			currentBinding.TCPPort, _ = strconv.Atoi(val)
		case "Password":
			currentBinding.Password = val
		case "BaudRate":
			currentBinding.SerialConf.BaudRate, _ = strconv.Atoi(val)
		case "DataBits":
			currentBinding.SerialConf.DataBits, _ = strconv.Atoi(val)
		case "Parity":
			currentBinding.SerialConf.Parity = val
		case "StopBits":
			currentBinding.SerialConf.StopBits = val
		case "FlowControl":
			currentBinding.SerialConf.FlowControl = val
		case "CharDelay":
			currentBinding.SerialConf.CharDelay, _ = strconv.Atoi(val)
		case "LineDelay":
			currentBinding.SerialConf.LineDelay, _ = strconv.Atoi(val)
		}
	}
	return nil
}

func (m *Manager) SaveConfig() error {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.saveConfigLocked()
}

func (m *Manager) saveConfigLocked() error {
	var sb strings.Builder
	for key, b := range m.Bindings {
		sb.WriteString(fmt.Sprintf("[%s]\n", key))
		sb.WriteString(fmt.Sprintf("SerialPort=%s\n", b.SerialPort))
		sb.WriteString(fmt.Sprintf("TCPPort=%d\n", b.TCPPort))
		sb.WriteString(fmt.Sprintf("Password=%s\n", b.Password))
		sb.WriteString(fmt.Sprintf("BaudRate=%d\n", b.SerialConf.BaudRate))
		sb.WriteString(fmt.Sprintf("DataBits=%d\n", b.SerialConf.DataBits))
		sb.WriteString(fmt.Sprintf("Parity=%s\n", b.SerialConf.Parity))
		sb.WriteString(fmt.Sprintf("StopBits=%s\n", b.SerialConf.StopBits))
		sb.WriteString(fmt.Sprintf("FlowControl=%s\n", b.SerialConf.FlowControl))
		sb.WriteString(fmt.Sprintf("CharDelay=%d\n", b.SerialConf.CharDelay))
		sb.WriteString(fmt.Sprintf("LineDelay=%d\n", b.SerialConf.LineDelay))
		sb.WriteString("\n")
	}
	return os.WriteFile(m.ConfigPath, []byte(sb.String()), 0644)
}

func (m *Manager) AddBinding(serialPort string, tcpPort int, password string, serialConf SerialConfig) error {
	if tcpPort < 1 || tcpPort > 65535 {
		return fmt.Errorf("invalid TCP port: must be between 1 and 65535")
	}

	key := fmt.Sprintf("%s:%d", serialPort, tcpPort)
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.Bindings[key]; ok {
		return fmt.Errorf("binding already exists")
	}

	m.Bindings[key] = &Binding{
		SerialPort: serialPort,
		TCPPort:    tcpPort,
		Password:   password,
		SerialConf: serialConf,
		Active:     false,
	}
	return m.saveConfigLocked()
}

func (m *Manager) RemoveBinding(key string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	if s, ok := m.Servers[key]; ok {
		s.Stop()
		delete(m.Servers, key)
	}
	delete(m.Bindings, key)
	return m.saveConfigLocked()
}

func (m *Manager) StartBinding(key string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	b, ok := m.Bindings[key]
	if !ok {
		return fmt.Errorf("binding not found")
	}

	if b.Active {
		return nil
	}

	server, err := NewSSHServer(b.TCPPort, b.SerialPort, b.Password, b.SerialConf)
	if err != nil {
		return err
	}

	if err := server.Start(m.hostKey); err != nil {
		return err
	}

	m.Servers[key] = server
	b.Active = true
	return nil
}

func (m *Manager) StopBinding(key string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	b, ok := m.Bindings[key]
	if !ok {
		return fmt.Errorf("binding not found")
	}

	if server, ok := m.Servers[key]; ok {
		server.Stop()
		delete(m.Servers, key)
	}
	b.Active = false
	return nil
}

func (m *Manager) GetBindings() map[string]*Binding {
	m.mu.RLock()
	defer m.mu.RUnlock()
	// Return a copy
	res := make(map[string]*Binding)
	for k, v := range m.Bindings {
		res[k] = v
	}
	return res
}
