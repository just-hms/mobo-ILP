package main

import (
	"fmt"

	"github.com/just-hms/mobo/pkg/opt"
)

func main() {
	res, maps := opt.Formalize([]*opt.Output{
		{
			Ones:      []uint{2, 3, 7, 12, 15},
			DontCares: []uint{4, 5, 13},
		},
		{
			Ones:      []uint{4, 7, 9, 11, 15},
			DontCares: []uint{6, 12, 14},
		},
	})

	fmt.Println(res)
	fmt.Println("-------------")
	fmt.Println(maps)
}
