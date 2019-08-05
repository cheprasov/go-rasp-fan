package main

import (
    "errors"
    "flag"
    "fmt"
    "github.com/stianeikeland/go-rpio"
    "io/ioutil"
    "log"
    "os"
    "os/exec"
    "regexp"
    "strconv"
    "strings"
    "time"
)

// Use pin 18, corresponds to physical pin 12 on the pi
const GPIOpin = 18

type TempConfig struct {
    temperature int
    workMs      int
    sleepMs     int
    isFanOn     bool
}

var tempConfigs = []TempConfig{
    {
        temperature: 50,
        workMs:      10000,
        sleepMs:     0,
    },
    {
        temperature: 40,
        workMs:      6000,
        sleepMs:     1000,
    },
    {
        temperature: 37,
        workMs:      2000,
        sleepMs:     1000,
    },
    {
        temperature: 35,
        workMs:      0,
        sleepMs:     3000,
    },
    {
        temperature: 20,
        workMs:      0,
        sleepMs:     5000,
    },
}

// temp=35.0'C
var regExpTemp = regexp.MustCompile(`temp\s*=\s*(\d+)(?:[,.]\d+)?`);

func temp() ([]byte, error) {
    return []byte("temp=35.0'C"), nil;
}

func getTemperature() (int, error) {
    out, err := exec.Command("/opt/vc/bin/vcgencmd", "measure_temp").Output()
    if err != nil {
        return 0, err;
    }

    matches := regExpTemp.FindStringSubmatch(string(out));
    if len(matches) != 2 {
        return 0, nil;
    }

    intValue, err := strconv.Atoi(matches[1]);
    if err != nil {
        return 0, err;
    }

    return intValue, nil;
}

var pin = rpio.Pin(GPIOpin)

func runFan() {
    pin.High()
}

func stopFan() {
    pin.Low()
}

func controlFan(t int) error {
    if t < 5 {
        return errors.New("strange temperature")
    }

    if len(tempConfigs) == 0 {
        return errors.New("empty temp config")
    }

    for _, cfg := range tempConfigs {
        if t >= cfg.temperature {
            fmt.Println("temp:", t, cfg);

            if cfg.workMs > 0 {
                runFan();
                time.Sleep(time.Duration(cfg.workMs) * time.Millisecond);
            }

            if cfg.sleepMs > 0 {
                stopFan();
                time.Sleep(time.Duration(cfg.sleepMs) * time.Millisecond);
            }

            return nil;
        }
    }

    return errors.New("can't find a temp config")
}

var pidFile string;

func init() {
    pidFilePointer := flag.String("pid", "", "Path to pid-file");
    flag.Parse();
    if *pidFilePointer == "" {
        println("Please provide all params for the script:");
        flag.PrintDefaults();
        log.Fatal();
    }

    pidFile = *pidFilePointer
}

func readFileContent(filename string) (string, error) {
    content, err := ioutil.ReadFile(filename)
    return string(content), err;
}

func isFileExists(filename string) bool {
    if _, err := os.Stat(filename); err != nil {
        if os.IsNotExist(err) {
            return false
        }
    }
    return true
}

var pid string;

func pidSaver(once bool) {
    for {
        now := time.Now().Unix();

        if isFileExists(pidFile) {
            content, err := readFileContent(pidFile)
            if err != nil {
                fmt.Println(err)
                time.Sleep(500 * time.Millisecond);
                continue;
            }

            if content != "" {
                data := strings.Split(content, ":");
                if len(data) == 2 && pid != data[1] {
                    lastTime, _ := strconv.ParseInt(data[0], 10, 64);
                    if now - lastTime < 6 {
                        os.Exit(0);
                        return;
                    }
                }
            }
        }

        strValue := fmt.Sprintf("%d:%s", now, pid);

        err := ioutil.WriteFile(pidFile, []byte(strValue), 0644)
        if err != nil {
            time.Sleep(500 * time.Millisecond)
            continue;
        }
        time.Sleep(5000 * time.Millisecond);

        if once {
            break;
        }
    }
}

func main() {
    pid = fmt.Sprintf("%d", os.Getpid());

    pidSaver(true);
    go pidSaver(false);

    if err := rpio.Open(); err != nil {
       fmt.Println(err)
       os.Exit(1)
    }
    // Unmap gpio memory when done
    defer rpio.Close()

    for {
        t, err := getTemperature();
        if err != nil || t == 0 {
            runFan();
            time.Sleep(5 * time.Second);
            continue;
        }
        if t > 70 {
            out, err := exec.Command("shutdown", "now").Output();
            if err != nil {
                log.Fatal("Can't shutdown system")
            }
            fmt.Println(out);
            return;
        }
        err = controlFan(t);
        if err != nil {
            runFan();
            time.Sleep(5 * time.Second);
            continue;
        }
    }
}
