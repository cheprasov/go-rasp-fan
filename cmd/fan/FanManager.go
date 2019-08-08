package fan

import (
    "fmt"
    "github.com/warthog618/gpio"
)

type FanManager struct {
    pin         *gpio.Pin
    runFanTemp  int
    stopFanTemp int
}

func CreateFanManager(GPIOPin uint8, runFanTemp, stopFanTemp int) (*FanManager, error) {
    if err := gpio.Open(); err != nil {
        return nil, err
    }

    fm := FanManager{
        pin:         gpio.NewPin(GPIOPin),
        runFanTemp:  runFanTemp,
        stopFanTemp: stopFanTemp,
    }
    fm.pin.Output()

    return &fm, nil
}

func (fm *FanManager) Close() error {
    return gpio.Close();
}

func (fm *FanManager) RunFan() {
    fm.pin.High()
}

func (fm *FanManager) StopFan() {
    fm.pin.Low()
}

func (fm *FanManager) ProcessTemp(t int) error {
    if t >= fm.runFanTemp {
        fm.RunFan();
        fmt.Println("temp: ", t, " => start fun on", fm.runFanTemp, "˚C")
    } else if t <= fm.stopFanTemp {
        fm.StopFan();
        fmt.Println("temp: ", t, " => stop fun on", fm.stopFanTemp, "˚C")
    }
    return nil
}
