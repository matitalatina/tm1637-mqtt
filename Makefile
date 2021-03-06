.PHONY: build build-arm

BUILD_CMD=go build -o dist/tm1637-mqtt cmd/sensor/main.go

build-arm6:
	GOOS=linux GOARCH=arm GOARM=6 $(BUILD_CMD)

build-arm:
	GOOS=linux GOARCH=arm GOARM=7 $(BUILD_CMD)

build:
	$(BUILD_CMD)

cp-mi:
	scp dist/tm1637-mqtt mi-rpi:tm1637-mqtt
	