#include <Servo.h>

const int HEAD_SERVO_PIN   = 12;
const int HEAD_ECHO_PIN    = 2; // Echo Pin
const int HEAD_TRIGGER_PIN = 3; // Trigger Pin
const int HEAD_POSITION_ZERO = 90;

Servo headServo;

inline int measureDistance() {
  digitalWrite(HEAD_TRIGGER_PIN, LOW); 
  delayMicroseconds(2); 

  digitalWrite(HEAD_TRIGGER_PIN, HIGH);
  delayMicroseconds(10); 
 
  digitalWrite(HEAD_TRIGGER_PIN, LOW);
  long duration = pulseIn(HEAD_ECHO_PIN, HIGH);
 
  //Calculate the distance (in cm) based on the speed of sound.
  return duration/58.2;
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
}

void loop() {
  int distance=measureDistance();
  Serial.println(distance);
  headServo.write(headDirection);
  delay(1000);
}
