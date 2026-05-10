package engine

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"os"
	"sync"

	"golang.org/x/crypto/ssh"
)

// Binding represents a mapping between a serial port and a TCP port
type Binding struct {
	SerialPort string `json:"serialPort"`
	TCPPort    int    `json:"tcpPort"`
	Password   string `json:"password"`
	Active     bool   `json:"active"`
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
	m.mu.Lock()
	defer m.mu.Unlock()
	return json.Unmarshal(data, &m.Bindings)
}

func (m *Manager) SaveConfig() error {
	m.mu.RLock()
	defer m.mu.RUnlock()
	data, err := json.MarshalIndent(m.Bindings, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(m.ConfigPath, data, 0644)
}

func (m *Manager) AddBinding(serialPort string, tcpPort int, password string) error {
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
		Active:     false,
	}
	return m.SaveConfig()
}

func (m *Manager) RemoveBinding(key string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	if s, ok := m.Servers[key]; ok {
		s.Stop()
		delete(m.Servers, key)
	}
	delete(m.Bindings, key)
	return m.SaveConfig()
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

	server, err := NewSSHServer(b.TCPPort, b.SerialPort, b.Password)
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
