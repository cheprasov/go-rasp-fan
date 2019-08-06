package config

import (
    "encoding/json"
    "io/ioutil"
)

type FanRule struct {
    Temp    int `json:"temp"`
    WorkMs  int `json:"workMs"`
    SleepMs int `json:"sleepMs"`
    Repeat  int `json:"repeat"`
}

type Config struct {
    GPIOPin              int       `json:"GPIOPin"`
    ShutdownTemp         int       `json:"shutdownTemp"`
    UpdateTempIntervalMs int       `json:"updateTempIntervalMs"`
    PidIntervalMs        int       `json:"pidIntervalMs"`
    FanRules             []FanRule `json:"fanRules"`
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
