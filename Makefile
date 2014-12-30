all: UPLOAD

clean:
	rm -rf .build

BUILD: src/robot.ino
	ino build

.build/uno-434b76f2/firmware.hex: BUILD

UPLOAD: .build/uno-434b76f2/firmware.hex
	ino upload
	ino serial

