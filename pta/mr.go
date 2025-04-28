package pta

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"time"
)

// ======= Implementação do Miller-Rabin =======

// MillerRabinTest verifica se um numero eh provavelmente primo
// usando o teste de primalidade de Miller-Rabin
// k eh o numero de iteracoes para aumentar a confiabilidade
func MillerRabinTest(n *big.Int, k int) bool {
	// Tratamento de casos especiais
	if n.Cmp(big.NewInt(2)) == 0 || n.Cmp(big.NewInt(3)) == 0 {
		return true
	}
	if n.Cmp(big.NewInt(2)) < 0 || new(big.Int).Mod(n, big.NewInt(2)).Cmp(big.NewInt(0)) == 0 {
		return false
	}

	// Escreve n-1 como 2^r * d onde d é ímpar
	r := 0
	d := new(big.Int).Sub(n, big.NewInt(1)) // d = n-1 inicialmente

	// Enquanto d eh par, dividir por 2
	for new(big.Int).Mod(d, big.NewInt(2)).Cmp(big.NewInt(0)) == 0 {
		d.Rsh(d, 1) // d = d/2
		r++
	}

	// Principal loop do Miller-Rabin
	for i := 0; i < k; i++ {
		if !millerRabinIteration(n, d, r) {
			return false // Definitivamente composto
		}
	}

	return true // Provavelmente primo
}

// millerRabinIteration realiza uma unica iteracao do teste
func millerRabinIteration(n, d *big.Int, r int) bool {
	// Escolhe um numero aleatorio a entre [2, n-2]
	nMinus2 := new(big.Int).Sub(n, big.NewInt(2))
	a, err := rand.Int(rand.Reader, nMinus2)
	if err != nil {
		// Fallback se rand.Int falhar
		t := time.Now().UnixNano()
		a = big.NewInt(t % nMinus2.Int64())
	}
	a.Add(a, big.NewInt(2)) // a esta agora entre 2 e n-2

	// Calcula x = a^d mod n
	x := new(big.Int).Exp(a, d, n)

	// Se x = 1 ou x = n-1, provavelmente eh primo
	one := big.NewInt(1)
	nMinus1 := new(big.Int).Sub(n, one)

	if x.Cmp(one) == 0 || x.Cmp(nMinus1) == 0 {
		return true
	}

	// Continua elevando ao quadrado x enquanto:
	// - r-1 > 0
	// - x != n-1
	// - x != 1
	for j := 0; j < r-1; j++ {
		// x = x^2 mod n
		x.Exp(x, big.NewInt(2), n)

		if x.Cmp(one) == 0 {
			// Encontramos uma raiz nao-trivial da unidade,
			// 	n é composto
			return false
		}

		if x.Cmp(nMinus1) == 0 {
			// Provavelmente primo
			return true
		}
	}

	// n eh composto
	return false
}

// GeneratePrimeNumber gera um numero primo com o tamanho de bits especificado
// usando o teste de Miller-Rabin
func GeneratePrimeNumber(bits int, candidato *big.Int) (*big.Int, int) {
	tentativas := 0
	for {
		tentativas++

		// Garantindo que o candidato tenha a quantidade de bits correto
		for candidato.BitLen() < bits {
			candidato.SetBit(candidato, bits-1, 1)
		}

		// Garantindo que o numero eh impar (um requisito para primos > 2)
		if candidato.Bit(0) == 0 {
			candidato.SetBit(candidato, 0, 1)
		}

		// Verificando se eh primo usando Miller-Rabin
		// O numero de iteracoes varia conforme o tamanho para aumentar a confiabilidade
		iteracoes := 20
		if bits > 256 {
			iteracoes = 30
		}
		if bits > 1024 {
			iteracoes = 40
		}

		if MillerRabinTest(candidato, iteracoes) {
			return candidato, tentativas
		}

		// Se nao for primo, incrementa por 2 e tentar novamente
		// Isto eh mais eficiente que gerar um novo numero aleatorio a cada tentativa
		// E como estamos usando o BBS e o LFG para gerar o candidato,
		// opto por nao "resetar" o gerador de numeros aleatorios
		candidato.Add(candidato, big.NewInt(2))
	}
}

// ======= Funcao chamada externamente =======
func MillerRabin(candidate *big.Int, bits int) *big.Int {
	fmt.Println("\nGerando número primo usando Miller-Rabin")
	fmt.Println("=======================================")

	// Medindo o execution time
	inicio := time.Now()

	// Gerando o primo
	prime, tentativas := GeneratePrimeNumber(bits, candidate)

	// Calculando o tempo decorrido
	duracao := time.Since(inicio)

	fmt.Printf("- Número gerado: %d bits\n", bits)
	fmt.Printf("- Tempo de execução: %s\n", duracao)
	fmt.Printf("- Tentativas: %d\n", tentativas)
	fmt.Printf("- Tamanho do número gerado: %d dígitos\n", len(prime.String()))

	// Verificando o tamanho real do numero gerado
	bitLength := prime.BitLen()

	fmt.Printf("- Tamanho real: %d bits\n", bitLength)
	fmt.Printf("- Valor decimal: %s\n", prime.String())

	binStr := fmt.Sprintf("%b", prime)
	fmt.Printf("- Binário: %s\n", binStr)

	return prime
}
