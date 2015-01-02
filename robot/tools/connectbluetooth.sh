#!/bin/sh

BT_ADDRESS=98:D3:31:B1:76:76
BT_DEVICE=0
BT_PIN=1234

# Add Serial Protcol
# sdptool add SP

# Connect ot th e device
rfcomm bind $BT_DEVICE $BT_ADDRESS

# Start the agent for the pin code
bluetooth-agent $BT_PIN &

# releasing the stuff
#rfcomm release $BT_ADDRESS


