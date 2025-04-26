package main

import (
	"PrimeNumGenerator/prng"
	"PrimeNumGenerator/pta"
)

func main() {
	bitSizes, generatedNumbers := prng.Lfg()
	for i, size := range bitSizes {
		pta.MillerRabin(generatedNumbers[i], size)
	}
	// prng.Bbs()
}
