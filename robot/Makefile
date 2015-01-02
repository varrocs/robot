FIRMWARE_FILE=.build/uno-434b76f2/firmware.hex 

all: UPLOAD 

clean:
	rm -rf .build

BUILD: src/robot.ino
	ino build

$(FIRMWARE_FILE): BUILD

UPLOAD: $(FIRMWARE_FILE)
	ino upload

