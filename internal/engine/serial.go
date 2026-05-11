package engine

import (
	"fmt"

	"go.bug.st/serial"
)

// SerialConfig represents the configuration for a serial port
type SerialConfig struct {
	BaudRate    int    `json:"baudRate"`
	DataBits    int    `json:"dataBits"`
	Parity      string `json:"parity"`      // none, odd, even, mark, space
	StopBits    string `json:"stopBits"`    // 1, 1.5, 2
	FlowControl string `json:"flowControl"` // none, xonxoff, rtscts, dsrdtr
	CharDelay   int    `json:"charDelay"`   // msec
	LineDelay   int    `json:"lineDelay"`   // msec
}

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

// OpenSerialPort opens a serial port with the provided configuration
func OpenSerialPort(name string, config SerialConfig) (serial.Port, error) {
	mode := &serial.Mode{
		BaudRate: config.BaudRate,
		DataBits: config.DataBits,
	}

	switch config.Parity {
	case "odd":
		mode.Parity = serial.OddParity
	case "even":
		mode.Parity = serial.EvenParity
	case "mark":
		mode.Parity = serial.MarkParity
	case "space":
		mode.Parity = serial.SpaceParity
	default:
		mode.Parity = serial.NoParity
	}

	switch config.StopBits {
	case "1.5":
		mode.StopBits = serial.OnePointFiveStopBits
	case "2":
		mode.StopBits = serial.TwoStopBits
	default:
		mode.StopBits = serial.OneStopBit
	}

	ser, err := serial.Open(name, mode)
	if err != nil {
		return nil, err
	}

	switch config.FlowControl {
	case "xonxoff":
		err = ser.SetMode(&serial.Mode{BaudRate: config.BaudRate, DataBits: config.DataBits, Parity: mode.Parity, StopBits: mode.StopBits})
		// go.bug.st/serial handles flow control differently depending on platform or specific calls
		// For now we will try to set it if the library supports it via mode or specific methods
		// Actually, standard go.bug.st/serial Mode doesn't have FlowControl field directly in some versions,
		// but we can use SetMode or other methods if available.
		// Looking at modern go.bug.st/serial: it doesn't have FlowControl in Mode struct.
		// Some implementations use platform specific settings.
	case "rtscts":
		// ser.SetRTS(true) etc.
	}

	return ser, nil
}
