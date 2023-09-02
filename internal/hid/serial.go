package hid

import (
	"github.com/xzzpig/simple-ip-kvm/internal/config"
	"go.bug.st/serial"
	"go.uber.org/zap"
)

var serialPort serial.Port
var dataQueue chan []byte
var closeQueue chan bool

func SetupVideoSerial() {
	mode := &serial.Mode{
		BaudRate: 19200,
	}
	port, err := serial.Open(config.GetConfig().SerialPort, mode)
	if err != nil {
		zap.L().Panic("open serial port error", zap.Error(err))
		return
	}
	serialPort = port
	dataQueue = make(chan []byte, 100)
	closeQueue = make(chan bool)

	go func() {
		for {
			select {
			case data := <-dataQueue:
				len, err := serialPort.Write(data)
				if err != nil {
					zap.L().Error("serial write error", zap.Error(err))
					continue
				}
				zap.L().Debug("serial write", zap.Binary("data", data), zap.Int("len", len))
			case <-closeQueue:
				return
			}
		}
	}()
}

func CLose() error {
	closeQueue <- true
	return serialPort.Close()
}

func Write(data []byte) {
	dataQueue <- data
}
