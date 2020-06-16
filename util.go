package main

import (
	"bufio"
	"io"
	"time"
)

func lineChannel(r io.Reader) (chan []byte, chan error) {
	lineChan := make(chan []byte, 100)
	errChan := make(chan error)

	reader := bufio.NewReader(r)
	go func() {
		for {
			line, err := reader.ReadBytes('\n')
			lineChan <- line

			if err != nil {
				errChan <- err
				break
			}
		}
		close(lineChan)
	}()

	return lineChan, errChan
}

func timerChannel(t time.Duration) chan bool {
	timerChan := make(chan bool)

	go func() {
		time.Sleep(t)
		timerChan <- true
	}()

	return timerChan
}
