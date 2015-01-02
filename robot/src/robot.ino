#include <Servo.h>
#include "message.h"

#define COMM 1

const int STATUS_LED       = 13;

const int HEAD_SERVO_PIN   = 12;
const int HEAD_ECHO_PIN    = 2; // Echo Pin
const int HEAD_TRIGGER_PIN = 3; // Trigger Pin
const int HEAD_POSITION_ZERO = 90;

const int MOTOR_ENA =  5;
const int MOTOR_ENB =  6;
const int MOTOR_IN1 =  8;
const int MOTOR_IN2 =  9;
const int MOTOR_IN3 = 10;
const int MOTOR_IN4 = 11;

const int DELAY_10CM = 650;

Servo headServo;

int measureDistance() {
  digitalWrite(HEAD_TRIGGER_PIN, LOW); 
  delayMicroseconds(2); 

  digitalWrite(HEAD_TRIGGER_PIN, HIGH);
  delayMicroseconds(10); 
 
  digitalWrite(HEAD_TRIGGER_PIN, LOW);
  long duration = pulseIn(HEAD_ECHO_PIN, HIGH);
 
  //Calculate the distance (in cm) based on the speed of sound.
  return duration/58.2;
}

void rightMotor(int speed) {
	// Stop
	if (speed==0) {
		digitalWrite(MOTOR_ENA, LOW);
		analogWrite(MOTOR_IN1, 0);
		analogWrite(MOTOR_IN2, 0);
	}
	else if (speed < 0) {
		digitalWrite(MOTOR_ENA, HIGH);
		analogWrite(MOTOR_IN1, -speed);
		analogWrite(MOTOR_IN2, 0);
	}
	else if (speed > 0) {
		digitalWrite(MOTOR_ENA, HIGH);
		analogWrite(MOTOR_IN1, 0);
		analogWrite(MOTOR_IN2, speed);
	}
}

void leftMotor(int speed) {
	// Stop
	if (speed==0) {
		digitalWrite(MOTOR_ENB, LOW);
		analogWrite(MOTOR_IN3, 0);
		analogWrite(MOTOR_IN4, 0);
	}
	else if (speed < 0) {
		digitalWrite(MOTOR_ENB, HIGH);
		analogWrite(MOTOR_IN3, -speed);
		analogWrite(MOTOR_IN4, 0);
	}
	else if (speed > 0) {
		digitalWrite(MOTOR_ENB, HIGH);
		analogWrite(MOTOR_IN3, 0);
		analogWrite(MOTOR_IN4, speed);
	}
}

int headDirection = 90;

void setup() {
  // Head setup
  //  Echo location
  pinMode(HEAD_TRIGGER_PIN, OUTPUT);
  pinMode(HEAD_ECHO_PIN, INPUT);
 
  //  Head servo
  headServo.attach(HEAD_SERVO_PIN);
  headServo.write(headDirection);

  // Motors
  pinMode(MOTOR_ENA, OUTPUT);
  pinMode(MOTOR_ENB, OUTPUT);

  pinMode(MOTOR_IN1, OUTPUT);
  pinMode(MOTOR_IN2, OUTPUT);
  pinMode(MOTOR_IN3, OUTPUT);
  pinMode(MOTOR_IN4, OUTPUT);

  rightMotor(0);
  leftMotor(0);

  // Communication
  pinMode(STATUS_LED, OUTPUT);
  Serial.begin(9600);
}

bool doNeedTurn(int direction) {
  headServo.write(direction);
  delay(500);
  int distance=measureDistance();

  return distance < 30;
}

void stop() {
	leftMotor(0);
	rightMotor(0);
}
void go() {
	leftMotor(255);
	rightMotor(255);
}

void simpleLoop() {
  stop();
  if (doNeedTurn(90)) {
  	bool needLeft = doNeedTurn(135);
	bool needRight = doNeedTurn(45);
	headServo.write(90);

	if (needLeft) 
	{
		rightMotor(-255);
	}
	if (needRight)
	{
		leftMotor(-255);
	}
	if (!needLeft && !needRight) {
		leftMotor(-255);
		rightMotor(-255);
	}
  }
  else {
  	go();
  }
  delay(DELAY_10CM);
}

ControlMessage receivedMessage;
ControlMessage sendMessage;
byte messageBuffer[5];

void doSendMessage() {
	Serial.write(sendMessage.len);
	Serial.write(sendMessage.type);
	Serial.write(sendMessage.param1);
	Serial.write('\n');
}

void sendDistance(int distance) {
	byte dist = (byte) min(distance, 255);
	sendMessage.len = 3;
	sendMessage.type = ECHO_RESPONSE;
	sendMessage.param1 = dist;
	doSendMessage();
}

void sendPong() {
	sendMessage.len = 3;
	sendMessage.type = PONG;
	sendMessage.param1 = 1;
	doSendMessage();
}

void executeReceivedMessage() {
	switch (receivedMessage.type) 
	{
	case ECHO_REQUEST:
		sendDistance(measureDistance());
		break;
	case SET_LEFT_MOTOR:
		leftMotor(receivedMessage.param1);
		break;
	case SET_RIGHT_MOTOR:
		rightMotor(receivedMessage.param1);
		break;
	case SET_HEAD:
		headServo.write(receivedMessage.param1);
		break;
	case WAIT:
		delay(receivedMessage.param1);
		break;
	case PING:
		sendPong();
		// fallthrough
	case PONG:
		// Flash the status led
		digitalWrite(STATUS_LED, HIGH);
		delay(100);
		digitalWrite(STATUS_LED, LOW);
	}
}

void commLoop() 
{
	digitalWrite(STATUS_LED, HIGH);
	delay(200);
	digitalWrite(STATUS_LED, LOW);

	ControlMessage message;
	if (Serial.readBytesUntil('\n', (char*)messageBuffer, sizeof messageBuffer))
	{
		int result = decode(messageBuffer, &receivedMessage);
		if (result == 0)
		{
			executeReceivedMessage();
		} else {
			Serial.print("ERROR\n");
			Serial.write(result);
			Serial.print("\n");
		}
	} 
}

void loop() 
{
#if COMM
commLoop();
#else
simpleLoop();
#endif
}

