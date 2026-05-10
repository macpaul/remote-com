package main

import (
	"context"
	"fmt"
	"log"
	"remote-com/internal/engine"
)

// App struct
type App struct {
	ctx     context.Context
	manager *engine.Manager
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	mgr, err := engine.NewManager("config.json")
	if err != nil {
		log.Printf("failed to initialize manager: %v", err)
	}
	a.manager = mgr
}

// ListPorts returns a list of available serial ports
func (a *App) ListPorts() ([]engine.SerialPortInfo, error) {
	return engine.ListSerialPorts()
}

// GetBindings returns all configured bindings
func (a *App) GetBindings() map[string]*engine.Binding {
	if a.manager == nil {
		return nil
	}
	return a.manager.GetBindings()
}

// AddBinding adds a new serial-to-tcp binding
func (a *App) AddBinding(serialPort string, tcpPort int, password string) error {
	if a.manager == nil {
		return fmt.Errorf("manager not initialized")
	}
	return a.manager.AddBinding(serialPort, tcpPort, password)
}

// RemoveBinding removes a binding
func (a *App) RemoveBinding(key string) error {
	if a.manager == nil {
		return fmt.Errorf("manager not initialized")
	}
	return a.manager.RemoveBinding(key)
}

// StartBinding starts the SSH server for a binding
func (a *App) StartBinding(key string) error {
	if a.manager == nil {
		return fmt.Errorf("manager not initialized")
	}
	return a.manager.StartBinding(key)
}

// StopBinding stops the SSH server for a binding
func (a *App) StopBinding(key string) error {
	if a.manager == nil {
		return fmt.Errorf("manager not initialized")
	}
	return a.manager.StopBinding(key)
}
