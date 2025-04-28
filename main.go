package main

import (
	"PrimeNumGenerator/prng"
	"PrimeNumGenerator/pta"
	"fmt"
	"os"
)

func LaggedFibonacci() {
	bitSizes, generatedNumbers := prng.Lfg()
	for i, size := range bitSizes {
		pta.MillerRabin(generatedNumbers[i], size)
		pta.Fermat(generatedNumbers[i], size)
	}
}

func Bbs() {
	bitSizes, generatedNumbers := prng.Bbs()
	for i, size := range bitSizes {
		pta.MillerRabin(generatedNumbers[i], size)
		pta.Fermat(generatedNumbers[i], size)
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Use: go run main.go [fibonacci|bbs]")
		return
	}

	switch os.Args[1] {
	case "fibonacci":
		LaggedFibonacci()
	case "bbs":
		Bbs()
	default:
		fmt.Println("Invalid option. Use: fibonacci, bbs")
		return
	}
}
