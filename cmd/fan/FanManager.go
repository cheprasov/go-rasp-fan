package fan

import (
    "../config"
    "../temp"
    "errors"
    "fmt"
    "github.com/warthog618/gpio"
    "time"
)

type FanManager struct {
    pin            *gpio.Pin
    fanRules       []config.FanRule
}

func CreateFanManager(GPIOPin uint8, fanRules []config.FanRule) (*FanManager, error) {
    if err := gpio.Open(); err != nil {
        return nil, err
    }

    fm := FanManager{
        pin:            gpio.NewPin(GPIOPin),
        fanRules:       fanRules,
    }
    fm.pin.Output()

    return &fm, nil
}

func (fm *FanManager) Close() error {
    // Unmap gpio memory when done
    return gpio.Close();
}

func (fm *FanManager) Run() {
    var t int
    var err error
    for {
        t, err = temp.GetTemperature()
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

            repeat := rule.Repeat
            if repeat < 1 {
                repeat = 1;
            }

            for repeat > 0 {
                if rule.RunMs > 0 {
                    fm.runFan()
                    time.Sleep(time.Duration(rule.RunMs) * time.Millisecond)
                }

                if rule.SleepMs > 0 {
                    fm.stopFan()
                    time.Sleep(time.Duration(rule.SleepMs) * time.Millisecond)
                }
                repeat--
            }

            return nil
        }
    }

    return errors.New("can't find a fan rule")
}
