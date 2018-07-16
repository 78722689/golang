package main

import "fmt"

import (
    "utility"
    "bufio"
    "os"
    "github.com/axgle/mahonia"
)

type Stmy struct {
    *utility.Stbase
}

func (m *Stmy) Run (id int) {
    fmt.Println("stmy::run id ", id)
}

func main() {
    fmt.Println("main")
    var v utility.Myer= &Stmy{&utility.Stbase{100}}

    utility.Runme(v)

    filename := "E:/Programing/golang/CFICCrawler/resource/configuration/funds.list"
    file, err:= os.Open(filename)
    if err != nil {
        fmt.Fprintf(os.Stderr, "WARN: Open file %s failed, %s\n", filename, err)
        return
    }
    defer file.Close()

    decoder := mahonia.NewDecoder("gbk")
    scanner := bufio.NewScanner(decoder.NewReader(file))
    for scanner.Scan() {
        fmt.Fprintf(os.Stdout, "%s\n", scanner.Text())
    }

}