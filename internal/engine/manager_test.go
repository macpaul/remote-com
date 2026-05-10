package engine

import (
	"os"
	"testing"
)

func TestManager_AddRemoveBinding(t *testing.T) {
	configPath := "test_config.json"
	defer os.Remove(configPath)
	defer os.Remove("host_key.pem")

	m, err := NewManager(configPath)
	if err != nil {
		t.Fatalf("Failed to create manager: %v", err)
	}

	err = m.AddBinding("COM1", 2222, "pass")
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
