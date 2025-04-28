// Esse arquivo traz a implementacao do Teste de Primalidade de Fermat
//  para verificar a primalidade e gerar numeros numero primos grandes.

package pta

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"time"
)

// FermatTest verifica se um numero eh provavelmente primo
// usando o teste de primalidade de Fermat.
// k eh o numero de iteracoes para aumentar a confiabilidade
func FermatTest(n *big.Int, k int) bool {
	// Tratamento de casos especiais
	if n.Cmp(big.NewInt(2)) == 0 || n.Cmp(big.NewInt(3)) == 0 {
		return true
	}
	if n.Cmp(big.NewInt(2)) < 0 || new(big.Int).Mod(n, big.NewInt(2)).Cmp(big.NewInt(0)) == 0 {
		return false
	}

	one := big.NewInt(1)
	nMinus1 := new(big.Int).Sub(n, one)

	for i := 0; i < k; i++ {
		nMinus2 := new(big.Int).Sub(n, big.NewInt(2))
		a, err := rand.Int(rand.Reader, nMinus2)
		if err != nil {
			t := time.Now().UnixNano()
			a = big.NewInt(t % nMinus2.Int64())
		}
		a.Add(a, big.NewInt(2)) // Garante que a >= 2

		// Calculamos a^(n-1) mod n
		result := new(big.Int).Exp(a, nMinus1, n)

		// Se o resultado != 1, entao definitivamente  eh composto
		if result.Cmp(one) != 0 {
			return false
		}
	}

	return true // Provavelmente primo
}

// GeneratePrimeNumberFermat gera um numero primo com o tamanho de bits especificado
// usando o Teste de Primalidade de Fermat.
func GeneratePrimeNumberFemart(bits int, candidato *big.Int) (*big.Int, int) {
	tentativas := 0
	for {
		tentativas++

		for candidato.BitLen() < bits {
			candidato.SetBit(candidato, bits-1, 1)
		}

		if candidato.Bit(0) == 0 {
			candidato.SetBit(candidato, 0, 1)
		}

		iteracoes := 20
		if bits > 256 {
			iteracoes = 30
		}
		if bits > 1024 {
			iteracoes = 40
		}

		if FermatTest(candidato, iteracoes) {
			return candidato, tentativas
		}

		candidato.Add(candidato, big.NewInt(2))
	}
}

// ======= Funcao chamada externamente =======
func Fermat(candidate *big.Int, bits int) *big.Int {
	fmt.Println("\nGerando número primo usando Fermat")
	fmt.Println("===================================")

	inicio := time.Now()

	prime, tentativas := GeneratePrimeNumberFemart(bits, candidate)

	duracao := time.Since(inicio)

	fmt.Printf("- Número gerado: %d bits\n", bits)
	fmt.Printf("- Tempo de execução: %s\n", duracao)
	fmt.Printf("- Tentativas: %d\n", tentativas)
	fmt.Printf("- Tamanho do número gerado: %d dígitos\n", len(prime.String()))

	bitLength := prime.BitLen()
	fmt.Printf("- Tamanho real: %d bits\n", bitLength)
	fmt.Printf("- Valor decimal: %s\n", prime.String())

	binStr := fmt.Sprintf("%b", prime)
	fmt.Printf("- Binário: %s\n", binStr)

	return prime
}
