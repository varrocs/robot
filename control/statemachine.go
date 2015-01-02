package main

import (
	"sync"
	"time"
)

const (
	DIRECTION_STRAIGHT = 90
	DIRECTION_LEFT     = 45
	DIRECTION_RIGHT    = 134

	FULL           = 127
	DISTANCE_LIMIT = 20
)

const (
	ROBOT10CM = 800 * time.Millisecond
)

type ListenerMap map[byte][]chan byte

type RobotState struct {
	sync.Mutex
	leftMotor     int
	rightMotor    int
	headDirection byte
	inputChannel  chan *ControlMessage
	outputChannel chan *ControlMessage
	distances     map[int]int // direction -> distance
	listeners     ListenerMap
}

func NewRobot(inputChannel chan *ControlMessage, outputChannel chan *ControlMessage) *RobotState {
	return &RobotState{leftMotor: 0, rightMotor: 0, headDirection: 90, inputChannel: inputChannel, outputChannel: outputChannel}
}

func (s *RobotState) Start() {
	s.Lock()
	defer s.Unlock()
	// Reset the robot
	s.setMotors(0, 0)
	s.turnHead(DIRECTION_STRAIGHT)
	// Start ECHO requests
	distanceChan := make(chan byte)
	s.addListener(MT_ECHO_RESPONSE, distanceChan)
	go func() {
		for {
			s.loop(distanceChan)
			time.Sleep(ROBOT10CM)
		}
	}()

}

func (s *RobotState) loop(distances <-chan byte) {
	s.Lock()
	defer s.Unlock()

	s.requestEcho()
	distance := <-distances

	if distance < DISTANCE_LIMIT {
		var turnLeft, turnRight bool
		s.setMotors(0, 0)

		s.turnHead(DIRECTION_LEFT)
		s.requestEcho()
		distance = <-distances
		turnRight = distance > DISTANCE_LIMIT

		s.turnHead(DIRECTION_RIGHT)
		s.requestEcho()
		distance = <-distances
		turnLeft = distance > DISTANCE_LIMIT

		if turnLeft {
			s.setMotors(-FULL, FULL)
		} else if turnRight {
			s.setMotors(FULL, -FULL)
		} else {
			s.setMotors(-FULL, -FULL)
		}
		time.Sleep(ROBOT10CM)
		s.setMotors(FULL, FULL)
	}

}

// -----------------------------------------------------------------

func (s *RobotState) addListener(messageType byte, callback chan byte) {
	s.Lock()
	defer s.Unlock()

	s.listeners[messageType] = append(s.listeners[messageType], callback)
}

func (s *RobotState) notifyListeners(messageType, value byte) {
	s.Lock()
	defer s.Unlock()

	for _, l := range s.listeners[messageType] {
		l <- value
	}
}

// ----------------------------------------------------------

func (s *RobotState) setMotor(isLeft bool, motorSpeed int) {
	s.Lock()
	defer s.Unlock()

	motorSpeedByte := byte(motorSpeed) + 127

	var messageType byte
	if isLeft {
		messageType = MT_SET_LEFT_MOTOR
		s.leftMotor = motorSpeed
	} else {
		messageType = MT_SET_RIGHT_MOTOR
		s.rightMotor = motorSpeed
	}

	msg := NewMessage(messageType, motorSpeedByte)
	s.outputChannel <- msg
}

func (s *RobotState) setMotors(left, right int) {
	s.Lock()
	defer s.Unlock()

	s.setMotor(true, left)
	s.setMotor(false, right)
}

func (s *RobotState) turnHead(direction byte) {
	s.Lock()
	defer s.Unlock()

	s.headDirection = direction
	msg := NewMessage(MT_SET_HEAD, direction)
	s.outputChannel <- msg
}

func (s *RobotState) requestEcho() {
	s.Lock()
	defer s.Unlock()

	msg := NewSimpleMessage(MT_ECHO_REQUEST)
	s.outputChannel <- msg
}
