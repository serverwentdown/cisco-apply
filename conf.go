package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"time"
)

type Device interface {
	ReadTimeout(time.Duration) ([]byte, error)
	WriteLine([]byte) error
	Close() error
}

var ErrorDeviceOperationNotSupported = errors.New("device operation not supported")

type ConfigurationLine interface {
	Apply(Device) error
}

type ConfigurationDirective struct {
	Name     string
	Argument string
}

type ConfigurationSimple []byte

func (d ConfigurationDirective) Apply(dev Device) error {
	if d.Name == "sleep" {
		seconds := 0
		_, err := fmt.Sscanf(string(d.Argument), "%d", &seconds)
		if err != nil {
			return fmt.Errorf("bad sleep arguments: %w", err)
		}
		time.Sleep(time.Duration(seconds) * time.Second)
	}
	return nil
}

func (s ConfigurationSimple) Apply(dev Device) error {
	err := dev.WriteLine([]byte(s))
	if err != nil {
		return err
	}
	response, err := dev.ReadTimeout(200 * time.Millisecond)
	if errors.Is(err, ErrorDeviceOperationNotSupported) {
	} else if err != nil {
		return err
	}
	fmt.Printf("%s", response)
	return nil
}

type ConfigurationReader bufio.Reader

func NewConfigurationReader(r io.Reader) *ConfigurationReader {
	return (*ConfigurationReader)(bufio.NewReader(r))
}

func (c *ConfigurationReader) Read() (ConfigurationLine, error) {
	line, err := (*bufio.Reader)(c).ReadBytes('\n')
	if err != nil {
		return nil, err
	}
	line = bytes.TrimSuffix(line, []byte{'\n'})

	if len(line) > 1 && line[0] == '!' && line[1] == ' ' {
		// Comment
		return ConfigurationSimple(string(line)), nil
	} else if len(line) > 0 && line[0] == '!' {
		// Directive
		line = line[1:]

		firstSep := bytes.IndexByte(line, ' ')
		if firstSep < 0 {
			return ConfigurationDirective{
				Name: string(line),
			}, nil
		}

		directive := line[:firstSep]
		rest := line[firstSep+1:]
		return ConfigurationDirective{
			Name:     string(directive),
			Argument: string(rest),
		}, nil
	} else {
		return ConfigurationSimple(line), nil
	}
}
