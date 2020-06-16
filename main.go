/*
apply interacts with a serial port to configure a Cisco product using files.

This tool does not attempt to understand the proprietary donkey format that is
the Cisco CLI. Instead, it provides some directives around it to make life
easier when crafting configuration files.

Basic usage is very simple: Pass a file containing Cisco commands and it will
write it out to serial port. It also prints the output back to stdout for
visual inspection of success.

apply also contains some additional features to help you in handling files. The
main feature is directives that perform certain actions. These are expressed
within a line beginning with "!".

Directives

	!sleep <seconds>

Sleep a fixed number of seconds.

	!assert <text>

Ensure that in the next 1 second, output from the serial port contains <text>.

*/
package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"time"
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s [-port PORT] [-writefile] FILE...\n", os.Args[0])
		flag.PrintDefaults()
	}
	port := flag.String("port", "/dev/ttyUSB0", "Serial port to write to")
	writefile := flag.Bool("writefile", false, "Treat PORT as file, and write to file")
	flag.Parse()
	confs := flag.Args()

	if len(confs) < 1 {
		flag.Usage()
		return
	}

	var device Device
	var err error
	if *writefile {
		device, err = NewDeviceFile(*port)
	} else {
		device, err = NewDeviceSerial(*port)
	}
	if err != nil {
		panic(err)
	}
	defer device.Close()

	for _, conf := range confs {
		file, err := os.Open(conf)
		if err != nil {
			panic(fmt.Errorf("file %s: %w", conf, err))
		}

		reader := NewConfigurationReader(file)
		lineNum := 0

		var line ConfigurationLine
		for {
			lineNum++

			line, err = reader.Read()
			if err != nil {
				break
			}
			if _, ok := line.(ConfigurationSimple); !ok {
				fmt.Printf("directive: %v\n", line)
			}

			err = line.Apply(device)
			if err != nil {
				break
			}
		}

		if err != nil && err != io.EOF {
			panic(fmt.Errorf("file %s line %d: %w", conf, lineNum, err))
		}

		// Attempt to read rest of output
		response, err := device.ReadTimeout(1000 * time.Millisecond)
		if errors.Is(err, ErrorDeviceOperationNotSupported) {
		} else if err != nil {
			panic(fmt.Errorf("device: %w", err))
		}
		fmt.Printf("%s", response)

		if err := file.Close(); err != nil {
			panic(fmt.Errorf("file %s: %w", conf, err))
		}
	}
}
