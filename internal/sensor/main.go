package sensor

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"periph.io/x/conn/v3/gpio/gpioreg"
	"periph.io/x/devices/v3/tm1637"
	"periph.io/x/host/v3"
)

type Config struct {
	Topic          string
	SensorPortPath string
	MqttBroker     string
	MqttUsername   string
	MqttPassword   string
}

func Start(c Config) {

	// Make sure periph is initialized.
	if _, err := host.Init(); err != nil {
		log.Fatal(err)
	}

	clk := gpioreg.ByName("GPIO3")
	data := gpioreg.ByName("GPIO2")
	if clk == nil || data == nil {
		log.Fatal("Failed to find pins")
	}
	dev, err := tm1637.New(clk, data)
	if err != nil {
		log.Fatalf("failed to initialize tm1637: %v", err)
	}
	if err := dev.SetBrightness(tm1637.Brightness10); err != nil {
		log.Fatalf("failed to set brightness on tm1637: %v", err)
	}

	opts := mqtt.NewClientOptions().AddBroker(c.MqttBroker).SetUsername(c.MqttUsername).SetPassword(c.MqttPassword).SetClientID("tm1637-mqtt")

	opts.AutoReconnect = true
	opts.SetKeepAlive(30 * time.Second)
	opts.SetPingTimeout(10 * time.Second)
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	var f mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
		fmt.Printf("MSG: %s\n", msg.Payload())
		strPayload, err := strconv.Atoi(string(msg.Payload()))
		if err != nil {
			log.Fatalf("unable to parse payload: %v", err)
		}
		if _, err := dev.Write(tm1637.Digits(splitToDigits(strPayload)...)); err != nil {
			log.Fatalf("failed to write to tm1637: %v", err)
		}
	}

	if token := client.Subscribe(c.Topic, 0, f); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	cmd := make(chan os.Signal, 1)
	signal.Notify(cmd, os.Interrupt, syscall.SIGTERM)
	<-cmd

}

func reverseInt(s []int) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}

func splitToDigits(n int) []int {
	if n > 9999 {
		return []int{9, 9, 9, 9}
	}
	if n <= 0 {
		return []int{0, 0, 0, 0}
	}
	var ret []int

	i := 0
	for i < 4 {
		var result = -1
		if n != 0 {
			result = n % 10
		}
		ret = append(ret, result)
		n /= 10
		i++
	}

	reverseInt(ret)

	return ret
}
