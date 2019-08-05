package main

import (
    "fmt"
    //"fmt"
    "log"
    "os/exec"

    //"github.com/stianeikeland/go-rpio"
    //"time"
)

// Use pin 18, corresponds to physical pin 12 on the pi
const GPIOpin = 18

//var pin = rpio.Pin(GPIOpin)

func main() {

    out, err := exec.Command("ls", "-la", "/Users/cheprasov/Projects/git/").Output()
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("The date is %s\n", out)

    return;
    //// Open and map memory to access gpio, check for errors
    //if err := rpio.Open(); err != nil {
    //    fmt.Println(err)
    //    os.Exit(1)
    //}
    //
    //// Unmap gpio memory when done
    //defer rpio.Close()
    //
    //// Set pin to output mode
    //pin.Output()
    //
    //// Toggle pin 20 times
    //for x := 0; x < 20; x++ {
    //    pin.Toggle()
    //    time.Sleep(time.Second / 5)
    //}
}