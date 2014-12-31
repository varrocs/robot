#include "message.h"
#include <string.h>

bool decode(const byte* raw_message, ControlMessage* message)
{
	// length
	size_t realLen = strlen((const char*)raw_message);
	byte len=raw_message[0];
	if (len != realLen || len < 2) {
		return false;
	}
	message->len = len;

	// type
	byte type = raw_message[1];
	if (type >= MESSAGE_TYPE_COUNT) {
		return false;
	}
	message->type =(ControlMessageType) type;

	if (len > 2) message->param1 = raw_message[2];
	if (len > 3) message->param2 = raw_message[3];
}
