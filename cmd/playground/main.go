package main

import "fmt"

func main() {
    s := make([]int, 0)
    for range 1000 {
       fmt.Println(len(s), cap(s)) // check the growth of 
       s = append(s, 1)
    }
}
