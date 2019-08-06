package main

import (
    "./cmd/config"
    "./cmd/fan"
    "flag"
    "io/ioutil"
    "log"
)

//var pidFile string
var cfg config.Config;

func init() {
    pidFilePointer := flag.String("pid", "", "Path to pid-file")
    configFilePointer := flag.String("config", "", "Path to config.json file")
    flag.Parse()
    if *pidFilePointer == "" {
        println("Please provide all params for the script:")
        flag.PrintDefaults()
        log.Fatal()
    }

    var err error;
    //pidFile = *pidFilePointer
    cfg, err = config.ReadConfig(*configFilePointer)
    if err != nil {
        log.Fatal(err)
    }
}

func readFileContent(filename string) (string, error) {
    content, err := ioutil.ReadFile(filename)
    return string(content), err
}

func main() {
    fanManager, err := fan.CreateFanManager(cfg.GPIOPin, cfg.FanRules)
    if err != nil {
        log.Fatal(err);
    }
    defer fanManager.Close()
    fanManager.Run();
}
