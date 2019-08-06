package fan

import (
    "errors"
    "fmt"
    "go-rasp-fan/cmd/config"
    "go-rasp-fan/cmd/temp"
    "time"
    gpio "github.com/stianeikeland/go-rpio"
)

type FanManager struct {
    pin gpio.Pin
    fanRules []config.FanRule
}

func CreateFanManager(pinId int, fanRules []config.FanRule) (*FanManager, error) {
    if err := gpio.Open(); err != nil {
        return nil, err
    }

    fm := FanManager{
        pin: gpio.Pin(pinId),
        fanRules: fanRules,
    }
    fm.pin.Output()

    return &fm, nil
}

func (fm *FanManager) Close() error {
    // Unmap gpio memory when done
    return gpio.Close();
}

func (fm *FanManager) Run() {
    for {
        t, err := temp.GetTemperature()
        if err != nil || t == 0 {
            fm.runFan()
            time.Sleep(5 * time.Second)
            continue
        }
        err = fm.processTemp(t)
        if err != nil {
            fm.runFan()
            time.Sleep(5 * time.Second)
            continue
        }
    }
}

func (fm *FanManager) runFan() {
    fm.pin.High()
}

func (fm *FanManager) stopFan() {
    fm.pin.Low()
}

func (fm *FanManager) processTemp(t int) error {
    if t < 5 {
        return errors.New("strange low temperature")
    }

    if len(fm.fanRules) == 0 {
        return errors.New("empty fun rules")
    }

    for _, rule := range fm.fanRules {
        if t >= rule.Temp {
            fmt.Println("temp:", t, rule)
            if rule.WorkMs > 0 {
                fm.runFan()
                time.Sleep(time.Duration(rule.WorkMs) * time.Millisecond)
            }

            if rule.SleepMs > 0 {
                fm.stopFan()
                time.Sleep(time.Duration(rule.SleepMs) * time.Millisecond)
            }

            return nil
        }
    }

    return errors.New("can't find a fan rule")
}