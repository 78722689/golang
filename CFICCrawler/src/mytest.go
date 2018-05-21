package main

//import "fmt"

func defer_test() *int{
    var y int=200
    y = y+100
    return &y
}

func pointerTest(v []int) {
    //fmt.Println(v[1])
    v[1] =100
    v=append(v,200)
}

type myST struct{
    x int
    //y *int
    u []int
}


func main() {
    var x = make([]int, 5, 5)
    x[1] = 10
    //fmt.Println(x)  // escaped to heap
    pointerTest(x)
    
    var st *myST
    st.x = 100
    //*st.y = 200
    st.u=x
    defer defer_test()
}