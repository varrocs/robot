#ifndef _MESSAGE_H_
#define _MESSAGE_H_

typedef unsigned char byte;

typedef enum ControlMessageType {
	ECHO_REQUEST = 1,
	ECHO_RESPONSE = 2,
	SET_LEFT_MOTOR = 3,
	SET_RIGHT_MOTOR = 4,
	SET_HEAD = 5,
	WAIT = 6,
	PING = 7,
	PONG = 8,
	MESSAGE_TYPE_COUNT
} ControlMessageType;

typedef struct ControlMessage {
	byte len;
	ControlMessageType type;
	byte param1;
	byte param2;

} ControlMessage;

int decode(const byte* raw_message, ControlMessage* message);
#endif
