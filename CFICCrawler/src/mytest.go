package main

import "fmt"

import "utility"

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
}