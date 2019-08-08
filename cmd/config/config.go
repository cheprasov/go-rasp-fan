package config

import (
    "encoding/json"
    "io/ioutil"
)

type Config struct {
    GPIOPin      uint8 `json:"GPIOPin"`
    ShutdownTemp int   `json:"shutdownTemp"`
    RunFanTemp   int   `json:"runFanTemp"`
    StopFanTemp  int   `json:"stopFanTemp"`
    WatchMs      int   `json:"watchMs"`
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
