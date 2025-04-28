package pta

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"time"
)

// ======= Implementação do Fermat =======

// FermatTest verifica se um número é provavelmente primo
// usando o teste de primalidade de Fermat
// k é o número de iterações para aumentar a confiabilidade
func FermatTest(n *big.Int, k int) bool {
	// Casos especiais
	if n.Cmp(big.NewInt(2)) == 0 || n.Cmp(big.NewInt(3)) == 0 {
		return true
	}
	if n.Cmp(big.NewInt(2)) < 0 || new(big.Int).Mod(n, big.NewInt(2)).Cmp(big.NewInt(0)) == 0 {
		return false
	}

	one := big.NewInt(1)
	nMinus1 := new(big.Int).Sub(n, one)

	for i := 0; i < k; i++ {
		// Escolher um número aleatório a entre [2, n-2]
		nMinus2 := new(big.Int).Sub(n, big.NewInt(2))
		a, err := rand.Int(rand.Reader, nMinus2)
		if err != nil {
			t := time.Now().UnixNano()
			a = big.NewInt(t % nMinus2.Int64())
		}
		a.Add(a, big.NewInt(2)) // Garantir que a >= 2

		// Calcular a^(n-1) mod n
		result := new(big.Int).Exp(a, nMinus1, n)

		// Se resultado != 1, então definitivamente composto
		if result.Cmp(one) != 0 {
			return false
		}
	}

	return true // Provavelmente primo
}

// ======= Funções para gerar números primos usando Fermat =======
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

// ======= Função principal para Fermat =======
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
