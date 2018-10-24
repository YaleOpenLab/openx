package main

import (
        "log"
        "github.com/tarm/serial"
        "bufio"
)

func main() {
        c := &serial.Config{
                Name: "/dev/ttyS0",
                Baud: 9600,
                Size: 8,
        }
        s, err := serial.OpenPort(c)
        if err != nil {
                log.Fatal(err)
        }
        n, err := s.Write([]byte("tests"))
        if err != nil {
                log.Fatal(err)
        }
        log.Print(n)
        for {
                r := bufio.NewReader(s)
                data, err := r.ReadBytes('\n')
                if err != nil {
                        log.Fatal(err)
                }
                log.Print(string(data[:]))
        }
}