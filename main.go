package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"

	"runtime/pprof"

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

var cpuprofile = flag.Bool("profile", false, "write cpu profile to \"data/cpu.prof\"")

func main() {
	flag.Parse()

	if *cpuprofile {
		f, err := os.Create("data/cpu.prof")
		if err != nil {
			log.Fatal(err)
		}

		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	outs := []*opt.Output{
		{Ones: []uint{1, 2, 3, 5}},
		{Ones: []uint{1, 5, 6, 7}},
	}
	log.Println("starting")
	res, maps := opt.Formalize(outs)

	fmt.Println(res)
	log.Println(maps)
}
