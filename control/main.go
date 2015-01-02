package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"sync"
	"time"

	"github.com/tarm/goserial"
)

const (
	PING_PONG_SECONDS     = 5
	DEFAULT_DEVICE        = "/dev/ttyACM0"
	MESSAGE_BUFFER_LENGTH = 4
)

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
		<-ch // Wait for the timer
		msg := NewSimpleMessage(MT_PING)
		msgSendCh <- msg
	}
}

func doControl(messageRecvCh chan *ControlMessage, messageSendCh chan *ControlMessage, robotChannel chan *ControlMessage) {
	for {
		msg := <-messageRecvCh
		typee := msg.MessageType
		if typee == MT_PONG {
			log.Println("PONG received!")
		} else if typee >= MT_INVALID {
			log.Println("Unknown message type! ", typee)
		} else {
			robotChannel <- msg
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
	messageRobotCh := make(chan *ControlMessage)

	var wg sync.WaitGroup

	robot := NewRobot(messageRobotCh, messageSendCh)

	startMessageIO(serial, messageRecvCh, messageSendCh)
	go doPingPong(messageSendCh)
	go doControl(messageRecvCh, messageSendCh, messageRobotCh)

	robot.Start()

	wg.Add(1)
	wg.Wait()
}
