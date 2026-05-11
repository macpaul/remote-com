package engine

import (
	"fmt"
	"os"
	"testing"
)

func TestManager_AddRemoveBinding(t *testing.T) {
	configPath := "test_settings.ini"
	defer os.Remove(configPath)
	defer os.Remove("host_key.pem")

	m, err := NewManager(configPath)
	if err != nil {
		t.Fatalf("Failed to create manager: %v", err)
	}

	conf := SerialConfig{BaudRate: 9600, DataBits: 8, Parity: "none", StopBits: "1"}
	err = m.AddBinding("COM1", 2222, "pass", conf)
	if err != nil {
		t.Fatalf("Failed to add binding: %v", err)
	}

	bindings := m.GetBindings()
	if len(bindings) != 1 {
		t.Errorf("Expected 1 binding, got %d", len(bindings))
	}

	key := "COM1:2222"
	if _, ok := bindings[key]; !ok {
		t.Errorf("Binding key %s not found", key)
	}

	err = m.RemoveBinding(key)
	if err != nil {
		t.Fatalf("Failed to remove binding: %v", err)
	}

	bindings = m.GetBindings()
	if len(bindings) != 0 {
		t.Errorf("Expected 0 bindings, got %d", len(bindings))
	}
}

func TestManager_InvalidPorts(t *testing.T) {
	configPath := "test_invalid_ports_settings.ini"
	defer os.Remove(configPath)
	defer os.Remove("host_key.pem")

	m, err := NewManager(configPath)
	if err != nil {
		t.Fatalf("Failed to create manager: %v", err)
	}

	invalidPorts := []int{-1, 0, 65536, 70000}
	conf := SerialConfig{BaudRate: 9600, DataBits: 8, Parity: "none", StopBits: "1"}
	for _, port := range invalidPorts {
		err := m.AddBinding("COM1", port, "pass", conf)
		if err == nil {
			t.Errorf("Expected error for port %d, but got nil", port)
		}
	}

	validPorts := []int{1, 80, 443, 1024, 65535}
	for _, port := range validPorts {
		err := m.AddBinding(fmt.Sprintf("COM_%d", port), port, "pass", conf)
		if err != nil {
			t.Errorf("Expected no error for port %d, but got %v", port, err)
		}
	}
}
