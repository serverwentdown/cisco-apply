package main

import (
	"time"

	"go.bug.st/serial"
)

type DeviceSerial struct {
	port    serial.Port
	lines   chan []byte
	lineErr chan error
}

func NewDeviceSerial(name string) (DeviceSerial, error) {
	serialMode := &serial.Mode{
		BaudRate: 9600,
		DataBits: 8,
		Parity:   serial.NoParity,
		StopBits: serial.OneStopBit,
	}
	port, err := serial.Open(name, serialMode)
	if err != nil {
		return DeviceSerial{}, err
	}
	lines, lineErr := lineChannel(port)
	return DeviceSerial{
		port:    port,
		lines:   lines,
		lineErr: lineErr,
	}, nil
}

func (s DeviceSerial) ReadTimeout(timeout time.Duration) ([]byte, error) {
	var err error
	buf := make([]byte, 0)
	t := timerChannel(timeout)
	for {
		// Read lines until timeout
		select {
		case line := <-s.lines:
			buf = append(buf, line...)
		case err = <-s.lineErr:
			return buf, err
		case <-t:
			return buf, err
		}
	}
	return nil, nil
}

func (s DeviceSerial) WriteLine(buf []byte) error {
	// TODO: handle n
	buf = append(buf, '\r')
	_, err := s.port.Write(buf)
	if err != nil {
		return err
	}
	return nil
}

func (s DeviceSerial) Close() error {
	return s.port.Close()
}
