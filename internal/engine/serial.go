package engine

import (
	"fmt"
	"go.bug.st/serial"
)

// SerialPortInfo represents basic information about a serial port
type SerialPortInfo struct {
	Name string `json:"name"`
}

// ListSerialPorts returns a list of available serial ports on the system
func ListSerialPorts() ([]SerialPortInfo, error) {
	ports, err := serial.GetPortsList()
	if err != nil {
		return nil, fmt.Errorf("failed to get serial ports: %w", err)
	}

	result := make([]SerialPortInfo, 0, len(ports))
	for _, port := range ports {
		result = append(result, SerialPortInfo{Name: port})
	}

	return result, nil
}

// OpenSerialPort opens a serial port with default settings (9600 8N1)
// We might want to make these configurable later.
func OpenSerialPort(name string) (serial.Port, error) {
	mode := &serial.Mode{
		BaudRate: 9600,
		DataBits: 8,
		Parity:   serial.NoParity,
		StopBits: serial.OneStopBit,
	}
	return serial.Open(name, mode)
}
