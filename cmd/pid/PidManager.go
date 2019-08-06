package pid

import (
    "fmt"
    "io/ioutil"
    "os"
    "strconv"
    "strings"
    "time"
)

func isFileExists(filename string) bool {
    if _, err := os.Stat(filename); err != nil {
        if os.IsNotExist(err) {
            return false
        }
    }
    return true
}

var pid = fmt.Sprintf("%d", os.Getpid())
var pidFile string

func readFileContent(filename string) (string, error) {
    content, err := ioutil.ReadFile(filename)
    return string(content), err
}

func pidSaver(once bool) {
    for {
        now := time.Now().Unix()
        if isFileExists(pidFile) {
            content, err := readFileContent(pidFile)
            if err != nil {
                fmt.Println(err)
                time.Sleep(500 * time.Millisecond)
                continue
            }

            if content != "" {
                data := strings.Split(content, ":")
                if len(data) == 2 && pid != data[1] {
                    lastTime, _ := strconv.ParseInt(data[0], 10, 64)
                    if now-lastTime < 6 {
                        os.Exit(0)
                        return
                    }
                }
            }
        }

        strValue := fmt.Sprintf("%d:%s", now, pid)
        err := ioutil.WriteFile(pidFile, []byte(strValue), 0644)
        if err != nil {
            time.Sleep(500 * time.Millisecond)
            continue
        }
        time.Sleep(5000 * time.Millisecond)
        if once {
            break
        }
    }
}