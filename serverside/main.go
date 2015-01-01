// Package main provides ...
package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"time"
)

type ControlMessage struct {
	Length      byte
	MessageType byte
	Param1      byte
}

const (
	PING_PONG_SECONDS = 10
	DEFAULT_DEVICE    = "/dev/ttyACM0"
	MESSAGE_LENGTH    = 3
)

const (
	MT_ECHO_REQUEST    = 0
	MT_ECHO_RESPONSE   = 1
	MT_SET_LEFT_MOTOR  = 2
	MT_SET_RIGHT_MOTOR = 3
	MT_SET_HEAD        = 4
	MT_WAIT            = 5
	MT_PING            = 6
	MT_PONG            = 7
)

func parseBuffer(buffer []byte) (*ControlMessage, error) {
	if len(buffer) < MESSAGE_LENGTH {
		return nil, errors.New("Too short buffer")
	}
	result := new(ControlMessage)
	result.Length = buffer[0]
	result.MessageType = buffer[1]
	result.Param1 = buffer[2]
	return result, nil

}

func serializeMessage(message *ControlMessage, result []byte) {
	result[0] = message.Length
	result[1] = message.MessageType
	result[2] = message.Param1
	result[3] = '\n'
}

func messageIO(device *os.File, msgRecvCh chan *ControlMessage, msgSendCh chan *ControlMessage) {
	// Start reading
	go func() {
		log.Println("Start reading")
		readBuffer := make([]byte, 4, 4)
		for {
			device.Read(readBuffer)
			message, err := parseBuffer(readBuffer)
			if err == nil {
				msgRecvCh <- message
			} else {
				log.Print(err)
			}
		}
	}()

	// Start sending
	go func() {
		log.Println("Start writing")
		writeBuffer := make([]byte, 4, 4)
		for {
			message := <-msgSendCh
			serializeMessage(message, writeBuffer)
			device.Write(writeBuffer)
		}
	}()
}

func startPingPong(msgSendCh chan *ControlMessage) {
	ch := time.NewTicker(PING_PONG_SECONDS * time.Second).C
	for {
		<-ch
		msg := ControlMessage{MESSAGE_LENGTH, MT_PING, 1}
		log.Println("Sending PING message", msg)
		msgSendCh <- &msg
	}
}

func startRead(readchan chan *ControlMessage) {
	for {
		message := <-readchan
		log.Println("Message received", message)
	}
}

func main() {
	fmt.Println("Starting comm")
	device, err := os.Open(DEFAULT_DEVICE)
	if err != nil {
		log.Fatal(err)
	}

	messageSendCh := make(chan *ControlMessage, 5)
	messageRecvCh := make(chan *ControlMessage, 5)

	messageIO(device, messageRecvCh, messageSendCh)
	go startPingPong(messageSendCh)
	go startRead(messageRecvCh)
}
