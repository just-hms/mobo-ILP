package main

import (
	"fmt"
	"log"
	"math/rand"

	"github.com/just-hms/mobo/pkg/opt"
	"golang.org/x/exp/maps"
)

func randomNonRepeatingNumbers(min, max, quantity int) []uint {
	if max-min+1 < quantity {
		panic("the range between min and max is too small to generate the required quantity of unique numbers")
	}

	generated := make(map[uint]bool)

	for len(generated) < quantity {
		num := uint(rand.Intn(max-min+1) + min)
		if !generated[num] {
			generated[num] = true
		}
	}

	return maps.Keys(generated)
}

func main() {
	outs := []*opt.Output{
		{Ones: randomNonRepeatingNumbers(0, 500, 300)},
		{Ones: randomNonRepeatingNumbers(0, 500, 300)},
		{Ones: randomNonRepeatingNumbers(0, 500, 300)},
		{Ones: randomNonRepeatingNumbers(0, 500, 300)},
		{Ones: randomNonRepeatingNumbers(0, 500, 300)},
		{Ones: randomNonRepeatingNumbers(0, 500, 300)},
		{Ones: randomNonRepeatingNumbers(0, 500, 300)},
		{Ones: randomNonRepeatingNumbers(0, 500, 300)},
		{Ones: randomNonRepeatingNumbers(0, 500, 300)},
		{Ones: randomNonRepeatingNumbers(0, 500, 300)},
		{Ones: randomNonRepeatingNumbers(0, 500, 300)},
		{Ones: randomNonRepeatingNumbers(0, 500, 300)},
		{Ones: randomNonRepeatingNumbers(0, 500, 300)},
		{Ones: randomNonRepeatingNumbers(0, 500, 300)},
		{Ones: randomNonRepeatingNumbers(0, 500, 300)},
		{Ones: randomNonRepeatingNumbers(0, 500, 300)},
		{Ones: randomNonRepeatingNumbers(0, 500, 300)},
		{Ones: randomNonRepeatingNumbers(0, 500, 300)},
		{Ones: randomNonRepeatingNumbers(0, 500, 300)},
		{Ones: randomNonRepeatingNumbers(0, 500, 300)},
		{Ones: randomNonRepeatingNumbers(0, 500, 300)},
		{Ones: randomNonRepeatingNumbers(0, 500, 300)},
		{Ones: randomNonRepeatingNumbers(0, 500, 300)},
		{Ones: randomNonRepeatingNumbers(0, 500, 300)},
		{Ones: randomNonRepeatingNumbers(0, 500, 300)},
		{Ones: randomNonRepeatingNumbers(0, 500, 300)},
		{Ones: randomNonRepeatingNumbers(0, 500, 300)},
		{Ones: randomNonRepeatingNumbers(0, 500, 300)},
		{Ones: randomNonRepeatingNumbers(0, 500, 300)},
		{Ones: randomNonRepeatingNumbers(0, 500, 300)},
		{Ones: randomNonRepeatingNumbers(0, 500, 300)},
		{Ones: randomNonRepeatingNumbers(0, 500, 300)},
		{Ones: randomNonRepeatingNumbers(0, 500, 300)},
		{Ones: randomNonRepeatingNumbers(0, 500, 300)},
		{Ones: randomNonRepeatingNumbers(0, 500, 300)},
		{Ones: randomNonRepeatingNumbers(0, 500, 300)},
		{Ones: randomNonRepeatingNumbers(0, 500, 300)},
		{Ones: randomNonRepeatingNumbers(0, 500, 300)},
		{Ones: randomNonRepeatingNumbers(0, 500, 300)},
		{Ones: randomNonRepeatingNumbers(0, 500, 300)},
		{Ones: randomNonRepeatingNumbers(0, 500, 300)},
		{Ones: randomNonRepeatingNumbers(0, 500, 300)},
		{Ones: randomNonRepeatingNumbers(0, 500, 300)},
		{Ones: randomNonRepeatingNumbers(0, 500, 300)},
		{Ones: randomNonRepeatingNumbers(0, 500, 300)},
		{Ones: randomNonRepeatingNumbers(0, 500, 300)},
		{Ones: randomNonRepeatingNumbers(0, 500, 300)},
		{Ones: randomNonRepeatingNumbers(0, 500, 300)},
		{Ones: randomNonRepeatingNumbers(0, 500, 300)},
		{Ones: randomNonRepeatingNumbers(0, 500, 300)},
		{Ones: randomNonRepeatingNumbers(0, 500, 300)},
	}
	log.Println("starting")
	res, maps := opt.Formalize(outs)

	fmt.Println(res)
	log.Println(maps)
}
