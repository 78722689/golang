package utility

import "fmt"

type Myer interface{
	Run(id int)
	Myfun()
}

type Stbase struct{
	X int
}

func (b * Stbase) Run (id int) {
	fmt.Println("stbase::run id", id)
}

func (b * Stbase)Myfun() {

}

func Runme (m Myer) {
	//fmt.Println("stbase::myfun")
	m.Run(20)
}