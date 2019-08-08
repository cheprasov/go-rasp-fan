package main

import (
    "./cmd/config"
    "./cmd/fan"
    "flag"
    "fmt"
    "go-rasp-fan/cmd/temp"
    "log"
    "os/exec"
    "time"
)

var cfg config.Config;

const (
    DefaultShutdownTemp = 70
    DefaultGPIOPin      = 18
    DefaultRunFanTemp   = 50
    DefaultStopFanTemp  = 36
)

func init() {
    configFilePointer := flag.String("config", "", "Path to config.json file")
    flag.Parse()

    var err error;
    cfg, err = config.ReadConfig(*configFilePointer)
    if err != nil {
        fmt.Println(err)
    }

    if cfg.GPIOPin == 0 {
        cfg.GPIOPin = DefaultGPIOPin
    }
    if cfg.RunFanTemp == 0 {
        cfg.RunFanTemp = DefaultRunFanTemp
    }
    if cfg.StopFanTemp == 0 {
        cfg.StopFanTemp = DefaultStopFanTemp
    }
    if cfg.ShutdownTemp == 0 {
        cfg.ShutdownTemp = DefaultShutdownTemp
    }
}

func shutdownNow() {
    output, err := exec.Command("shutdown", "now").Output()
    if err != nil {
        log.Println(err)
    }
    fmt.Println(string(output))
}

func main() {
    fanManager, err := fan.CreateFanManager(cfg.GPIOPin, cfg.RunFanTemp, cfg.StopFanTemp)
    if err != nil {
        log.Fatal(err);
    }
    defer fanManager.Close()

    for {
        t, err := temp.GetTemperature()

        if t >= cfg.ShutdownTemp {
            shutdownNow();
        }

        if err != nil || t == 0 {
            // Can't get temp
            fanManager.RunFan()
            time.Sleep(time.Duration(cfg.WatchMs) * time.Millisecond)
            continue
        }

        err = fanManager.ProcessTemp(t)
        if err != nil {
            // Can't process temp
            fanManager.RunFan();
            time.Sleep(time.Duration(cfg.WatchMs) * time.Millisecond)
            continue
        }
    }
}
