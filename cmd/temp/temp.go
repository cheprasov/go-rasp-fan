package temp

import (
    "errors"
    "os/exec"
    "regexp"
    "strconv"
)

var regExpTemp = regexp.MustCompile(`temp\s*=\s*(\d+)(?:[,.]\d+)?`)

func GetTemperature() (int, error) {
    out, err := exec.Command("/opt/vc/bin/vcgencmd", "measure_temp").Output()
    if err != nil {
        return 0, err
    }

    matches := regExpTemp.FindStringSubmatch(string(out))
    if len(matches) != 2 {
        return 0, errors.New("wrong matches count, can't parse")
    }

    intValue, err := strconv.Atoi(matches[1])
    if err != nil {
        return 0, err
    }

    return intValue, nil
}