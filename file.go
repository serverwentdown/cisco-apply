package main

import (
	"fmt"
	"os"
	"time"
)

type DeviceFile struct {
	file *os.File
}

func NewDeviceFile(name string) (DeviceFile, error) {
	file, err := os.OpenFile(name, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return DeviceFile{}, err
	}
	return DeviceFile{
		file: file,
	}, nil
}

func (f DeviceFile) ReadTimeout(timeout time.Duration) ([]byte, error) {
	return nil, fmt.Errorf("%w: Read", ErrorDeviceOperationNotSupported)
}

func (f DeviceFile) WriteLine(buf []byte) error {
	buf = append(buf, '\n')
	f.file.Write(buf)
	return nil
}

func (f DeviceFile) Close() error {
	return f.file.Close()
}
