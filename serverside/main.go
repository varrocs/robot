package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"sync"
	"time"

	"github.com/tarm/goserial"
)

type ControlMessage struct {
	Length      byte
	MessageType byte
	Param1      byte
}

const (
	PING_PONG_SECONDS     = 5
	DEFAULT_DEVICE        = "/dev/ttyACM0"
	MESSAGE_BUFFER_LENGTH = 4
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
	if len(buffer) < MESSAGE_BUFFER_LENGTH {
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

func startMessageIO(device io.ReadWriteCloser, msgRecvCh chan *ControlMessage, msgSendCh chan *ControlMessage) {
	// Start reading
	go func() {
		log.Println("Start reading")
		var reader *bufio.Reader = bufio.NewReader(device)
		for {
			readBuffer, err := reader.ReadSlice('\n')
			if err == io.EOF {
				time.Sleep(1 * time.Second)
				continue
			} else if err != nil {
				log.Println(err)
				continue
			}
			log.Println(" >>>>>>>>>>> RAW MESSAGE: ", readBuffer)

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
		writeBuffer := make([]byte, MESSAGE_BUFFER_LENGTH, MESSAGE_BUFFER_LENGTH)
		for {
			message := <-msgSendCh
			serializeMessage(message, writeBuffer)
			log.Println(" <<<<<<<<<<< RAW MESSAGE: ", writeBuffer)
			_, err := device.Write(writeBuffer)
			if err != nil {
				log.Println(err)
			}
		}
	}()
}

func doPingPong(msgSendCh chan *ControlMessage) {
	ch := time.NewTicker(PING_PONG_SECONDS * time.Second).C
	for {
		<-ch
		msg := ControlMessage{MESSAGE_BUFFER_LENGTH - 1, MT_PING, 1}
		msgSendCh <- &msg
	}
}

func doRead(readchan chan *ControlMessage) {
	for {
		message := <-readchan
		switch message.MessageType {
		case MT_PONG:
			log.Println("PONG received!")
		}
	}
}

func main() {
	deviceName := flag.String("serialdevice", DEFAULT_DEVICE, "The comm with the robot will happen through this device")
	flag.Parse()

	fmt.Println("Starting comm")
	serialConfig := &serial.Config{Name: *deviceName, Baud: 9600}
	serial, err := serial.OpenPort(serialConfig)
	if err != nil {
		log.Fatal(err)
	}

	messageSendCh := make(chan *ControlMessage)
	messageRecvCh := make(chan *ControlMessage)

	var wg sync.WaitGroup

	startMessageIO(serial, messageRecvCh, messageSendCh)
	go doPingPong(messageSendCh)
	go doRead(messageRecvCh)

	wg.Add(1)
	wg.Wait()
}
