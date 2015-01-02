package main

import (
	"errors"
)

const (
	MT_ECHO_REQUEST    = 1
	MT_ECHO_RESPONSE   = 2
	MT_SET_LEFT_MOTOR  = 3
	MT_SET_RIGHT_MOTOR = 4
	MT_SET_HEAD        = 5
	MT_WAIT            = 6
	MT_PING            = 7
	MT_PONG            = 8

	MT_INVALID = 9
)

type ControlMessage struct {
	Length      byte
	MessageType byte
	Param1      byte
}

func NewSimpleMessage(messageType byte) *ControlMessage {
	return &ControlMessage{Length: 3, MessageType: messageType, Param1: 1}
}

func NewMessage(messageType byte, param byte) *ControlMessage {
	return &ControlMessage{Length: 3, MessageType: messageType, Param1: param}
}

func parseBuffer(buffer []byte) (*ControlMessage, error) {
	if len(buffer) < MESSAGE_BUFFER_LENGTH {
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
