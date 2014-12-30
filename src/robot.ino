#include <Servo.h>

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
	Serial.print("right - ");
	Serial.println(speed);
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
	Serial.print("left - ");
	Serial.println(speed);

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
  Serial.begin(9600);
  
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
}

bool doNeedTurn(int direction) {
  headServo.write(direction);
  delay(500);
  int distance=measureDistance();
  Serial.print(direction);
  Serial.print(" - ");
  Serial.println(distance);

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

void loop() {
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
  Serial.println();
}
