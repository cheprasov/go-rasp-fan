package config

import (
    "encoding/json"
    "io/ioutil"
)

type FanRule struct {
    Temp    int `json:"temp"`
    RunMs   int `json:"runMs"`
    SleepMs int `json:"sleepMs"`
    Repeat  int `json:"repeat"`
}

type Config struct {
    GPIOPin       uint8     `json:"GPIOPin"`
    ShutdownTemp  int       `json:"shutdownTemp"`
    PidIntervalMs int       `json:"pidIntervalMs"`
    RunTemp       int       `json:"runTemp"`
    StopTemp      int       `json:"stopTemp"`
    FanRules      []FanRule `json:"fanRules"`
}

func ReadConfig(filename string) (Config, error) {
    var config Config;
    jsonData, err := ioutil.ReadFile(filename)
    if err != nil {
        return config, err;
    }

    err = json.Unmarshal(jsonData, &config);
    if err != nil {
        return config, err
    }

    return config, nil;
}
