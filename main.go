package main

import (
	"fmt"

	"github.com/just-hms/mobo/pkg/opt"
)

func main() {
	res, maps := opt.Formalize([]*opt.Output{
		{Ones: []uint{1, 2, 3, 5}},
		{Ones: []uint{1, 5, 6, 7}},
	})

	fmt.Println(res)
	fmt.Println("\\", maps)
}
