package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"sync"
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
	MESSAGE_LENGTH    = 4
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
		return nil, errors.New("Too show buffer")
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

func startMessageIO(device *os.File, msgRecvCh chan *ControlMessage, msgSendCh chan *ControlMessage) {
	// Start reading
	go func() {
		log.Println("Start reading")
		readBuffer := make([]byte, MESSAGE_LENGTH, MESSAGE_LENGTH)
		for {
			length, err := device.Read(readBuffer)
			if err != nil {
				log.Println(err)
				return
			}

			if length == 0 {
				log.Println("Zero message read, finishing read cycle")
				return
			}

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
		writeBuffer := make([]byte, MESSAGE_LENGTH, MESSAGE_LENGTH)
		for {
			message := <-msgSendCh
			serializeMessage(message, writeBuffer)
			device.Write(writeBuffer)
		}
	}()
}

func doPingPong(msgSendCh chan *ControlMessage) {
	ch := time.NewTicker(PING_PONG_SECONDS * time.Second).C
	for {
		<-ch
		msg := ControlMessage{MESSAGE_LENGTH, MT_PING, 1}
		log.Println("Sending PING message", msg)
		msgSendCh <- &msg
	}
}

func doRead(readchan chan *ControlMessage) {
	for {
		message := <-readchan
		log.Println("Message received", message)
	}
}

func main() {
	deviceName := flag.String("serialdevice", DEFAULT_DEVICE, "The comm with the robot will happen through this device")
	flag.Parse()

	fmt.Println("Starting comm")
	device, err := os.Open(*deviceName)
	if err != nil {
		log.Fatal(err)
	}

	messageSendCh := make(chan *ControlMessage)
	messageRecvCh := make(chan *ControlMessage)

	var wg sync.WaitGroup
	startMessageIO(device, messageRecvCh, messageSendCh)
	go doPingPong(messageSendCh)
	go doRead(messageRecvCh)
	wg.Add(1)
	wg.Wait()
}
