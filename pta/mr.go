package pta

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"time"
)

// ======= Implementação do Miller-Rabin =======

// MillerRabinTest verifica se um número é provavelmente primo
// usando o teste de primalidade de Miller-Rabin
// k é o número de iterações para aumentar a confiabilidade
func MillerRabinTest(n *big.Int, k int) bool {
	// Casos especiais
	if n.Cmp(big.NewInt(2)) == 0 || n.Cmp(big.NewInt(3)) == 0 {
		return true
	}
	if n.Cmp(big.NewInt(2)) < 0 || new(big.Int).Mod(n, big.NewInt(2)).Cmp(big.NewInt(0)) == 0 {
		return false
	}

	// Escrever n-1 como 2^r * d onde d é ímpar
	r := 0
	d := new(big.Int).Sub(n, big.NewInt(1)) // d = n-1 inicialmente

	// Enquanto d é par, dividir por 2
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

// millerRabinIteration realiza uma única iteração do teste
func millerRabinIteration(n, d *big.Int, r int) bool {
	// Escolher um número aleatório a entre [2, n-2]
	nMinus2 := new(big.Int).Sub(n, big.NewInt(2))
	a, err := rand.Int(rand.Reader, nMinus2)
	if err != nil {
		// Fallback se rand.Int falhar
		t := time.Now().UnixNano()
		a = big.NewInt(t % nMinus2.Int64())
	}
	a.Add(a, big.NewInt(2)) // a está agora entre 2 e n-2

	// Calcular x = a^d mod n
	x := new(big.Int).Exp(a, d, n)

	// Se x = 1 ou x = n-1, provavelmente é primo
	one := big.NewInt(1)
	nMinus1 := new(big.Int).Sub(n, one)

	if x.Cmp(one) == 0 || x.Cmp(nMinus1) == 0 {
		return true
	}

	// Continuar elevando ao quadrado x enquanto:
	// - r-1 > 0
	// - x != n-1
	// - x != 1
	for j := 0; j < r-1; j++ {
		// x = x^2 mod n
		x.Exp(x, big.NewInt(2), n)

		if x.Cmp(one) == 0 {
			// Encontramos uma raiz não-trivial da unidade, n é composto
			return false
		}

		if x.Cmp(nMinus1) == 0 {
			// Provavelmente primo
			return true
		}
	}

	// n é composto
	return false
}

// ======= Funções para gerar números primos =======

// GerarNumeroPrimo gera um número primo com o tamanho de bits especificado
// usando o teste de Miller-Rabin
func GerarNumeroPrimo(bits int, candidato *big.Int) (*big.Int, int) {
	tentativas := 0
	for {
		tentativas++

		// Garantir que o candidato tenha o número de bits correto
		for candidato.BitLen() < bits {
			candidato.SetBit(candidato, bits-1, 1)
		}

		// Garantir que o número é ímpar (um requisito para primos > 2)
		if candidato.Bit(0) == 0 {
			candidato.SetBit(candidato, 0, 1)
		}

		// Verificar se é primo usando Miller-Rabin
		// O número de iterações varia conforme o tamanho para equilibrar segurança e desempenho
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

		// Se não for primo, incrementar por 2 e tentar novamente
		// Isto é mais eficiente que gerar um novo número aleatório a cada tentativa
		candidato.Add(candidato, big.NewInt(2))
	}
}

// ======= Função principal =======

func MillerRabin(candidate *big.Int, bits int) *big.Int {
	fmt.Println("Gerando número primo usando Miller-Rabin")
	fmt.Println("=======================================")

	// Cabeçalho da tabela
	fmt.Printf("| %-6s | %-10s | %-15s | %-20s | %-15s |\n",
		"Bits", "Método", "Tentativas", "Tempo (segundos)", "Dígitos Decimais")
	fmt.Println("|--------|------------|-----------------|----------------------|-----------------|")

	// Medir tempo de execução
	inicio := time.Now()

	// Gerar primo
	primo, tentativas := GerarNumeroPrimo(bits, candidate)

	// Calcular tempo decorrido
	duracao := time.Since(inicio)

	// Mostrar resultados
	fmt.Printf("| %-6d | %-10s | %-15d | %-20.4f | %-15d |\n",
		bits, "Miller-Rabin", tentativas, duracao.Seconds(), len(primo.String()))

	// Verificar o tamanho real do número gerado
	bitLength := primo.BitLen()

	// Mostrar detalhes do primo gerado
	fmt.Printf("|- Tamanho real: %d bits\n", bitLength)
	fmt.Printf("|- Valor decimal: %s\n", primo.String())

	// Mostrar parte inicial e final da representação binária
	binStr := fmt.Sprintf("%b", primo)
	if len(binStr) > 64 {
		fmt.Printf("|- Binário: %s...%s\n", binStr[:32], binStr[len(binStr)-32:])
	} else {
		fmt.Printf("|- Binário: %s\n", binStr)
	}
	fmt.Println("|--------|------------|-----------------|----------------------|-----------------|")

	return primo
}
