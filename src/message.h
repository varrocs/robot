#ifndef _MESSAGE_H_
#define _MESSAGE_H_

typedef unsigned char byte;

typedef enum ControlMessageType {
	ECHO_REQUEST = 0,
	ECHO_RESPONSE = 1,
	SET_LEFT_MOTOR = 2,
	SET_RIGHT_MOTOR = 3,
	SET_HEAD = 4,
	WAIT = 5,
	PING = 6,
	PONG = 7,
	MESSAGE_TYPE_COUNT
} ControlMessageType;

typedef struct ControlMessage {
	byte len;
	ControlMessageType type;
	byte param1;
	byte param2;

} ControlMessage;

bool decode(const byte* raw_message, ControlMessage* message);
#endif
