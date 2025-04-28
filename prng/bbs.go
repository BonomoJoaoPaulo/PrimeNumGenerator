package prng

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"time"
)

// BlumBlumShub implementa o algoritmo BBS para gerar números pseudoaleatórios
type BlumBlumShub struct {
	p, q    *big.Int // Primos congruentes a 3 mod 4
	n       *big.Int // n = p * q
	state   *big.Int // Estado atual x_i
	bitSize int      // Tamanho desejado em bits
}

// NewBBS cria um novo gerador BBS
func NewBBS(bitSize int) *BlumBlumShub {
	// Calcular quantos bits cada primo deve ter (aproximadamente metade do tamanho total)
	primeBits := (bitSize + 1) / 2

	// Gerar primos p e q, ambos congruentes a 3 mod 4
	p := generateSafePrime(primeBits)
	q := generateSafePrime(primeBits)

	// Garantir que p != q
	for p.Cmp(q) == 0 {
		q = generateSafePrime(primeBits)
	}

	// Calcular n = p * q
	n := new(big.Int).Mul(p, q)

	// Gerar um valor inicial (semente) x_0 que seja coprimo com n
	seed := generateSeed(n)

	bbs := &BlumBlumShub{
		p:       p,
		q:       q,
		n:       n,
		state:   seed,
		bitSize: bitSize,
	}

	return bbs
}

// generateSafePrime gera um número primo p tal que p ≡ 3 (mod 4)
func generateSafePrime(bits int) *big.Int {
	three := big.NewInt(3)
	four := big.NewInt(4)

	for {
		// Gerar um número primo aleatório com o tamanho especificado
		p, err := rand.Prime(rand.Reader, bits)
		if err != nil {
			// Fallback se rand.Prime falhar
			p = generateFallbackPrime(bits)
		}

		// Verificar se p ≡ 3 (mod 4)
		if new(big.Int).Mod(p, four).Cmp(three) == 0 {
			return p
		}
	}
}

// generateFallbackPrime gera um número primo quando rand.Prime falha
// Nota: Este é um método simplificado e não deve ser usado em produção
func generateFallbackPrime(bits int) *big.Int {
	// Iniciar com um número ímpar aleatório
	candidate := big.NewInt(0)
	max := new(big.Int).Lsh(big.NewInt(1), uint(bits))

	for {
		// Gerar um número aleatório
		t := time.Now().UnixNano()
		candidate.SetInt64(t)
		candidate.Mod(candidate, max)

		// Garantir que é ímpar
		if candidate.Bit(0) == 0 {
			candidate.Add(candidate, big.NewInt(1))
		}

		// Garantir que é congruente a 3 mod 4
		if new(big.Int).Mod(candidate, big.NewInt(4)).Cmp(big.NewInt(3)) != 0 {
			candidate.Add(candidate, big.NewInt(2))
		}

		// Verificar se é provavelmente primo
		if candidate.ProbablyPrime(20) {
			return candidate
		}

		// Tentar o próximo número congruente a 3 mod 4
		candidate.Add(candidate, big.NewInt(4))
		time.Sleep(time.Nanosecond) // Variar o timestamp
	}
}

// generateSeed gera um valor inicial x_0 que seja coprimo com n
func generateSeed(n *big.Int) *big.Int {
	one := big.NewInt(1)

	for {
		// Gerar um número aleatório entre 2 e n-1
		seed, err := rand.Int(rand.Reader, new(big.Int).Sub(n, big.NewInt(2)))
		if err != nil {
			// Fallback se rand.Int falhar
			t := time.Now().UnixNano()
			seed = big.NewInt(t)
			seed.Mod(seed, n)
		}

		seed.Add(seed, big.NewInt(2)) // Agora seed está entre 2 e n-1

		// Verificar se o seed é coprimo com n usando GCD
		gcd := new(big.Int).GCD(nil, nil, seed, n)

		if gcd.Cmp(one) == 0 {
			// Calcular x_0 = seed^2 mod n para iniciar a sequência
			x0 := new(big.Int).Exp(seed, big.NewInt(2), n)
			return x0
		}
	}
}

// NextState calcula o próximo estado x_(i+1) = x_i^2 mod n
func (bbs *BlumBlumShub) NextState() *big.Int {
	// x_(i+1) = x_i^2 mod n
	bbs.state = new(big.Int).Exp(bbs.state, big.NewInt(2), bbs.n)
	return new(big.Int).Set(bbs.state)
}

// NextBit gera o próximo bit (o bit de paridade do estado)
func (bbs *BlumBlumShub) NextBit() uint {
	// Atualizar o estado
	bbs.NextState()

	// Retornar o bit de paridade (LSB)
	return bbs.state.Bit(0)
}

// Next gera um número pseudoaleatório com o tamanho aproximado de bitSize
func (bbs *BlumBlumShub) Next() *big.Int {
	result := big.NewInt(0)

	// Gerar bitSize bits para formar o número
	for i := 0; i < bbs.bitSize; i++ {
		bit := bbs.NextBit()

		// Deslocar o resultado e adicionar o novo bit
		result.Lsh(result, 1)
		if bit == 1 {
			result.Or(result, big.NewInt(1))
		}
	}

	return result
}

func Bbs() ([]int, []*big.Int) {
	// Tamanhos de bits para testar
	bitSizes := []int{40, 56, 80, 128, 168, 224, 256, 512, 1024, 2048, 4096}
	generatedNumbers := make([]*big.Int, len(bitSizes))

	fmt.Println("Gerando números pseudoaleatórios com Blum Blum Shub")
	fmt.Println("=================================================")

	for i, bits := range bitSizes {
		fmt.Printf("\nGerando número de %d bits:\n", bits)

		// Criar um novo gerador para cada tamanho de bits
		fmt.Printf("- Gerando primos p e q (isso pode levar alguns instantes)...\n")
		init_time := time.Now()
		bbs := NewBBS(bits)
		elapsed_time := time.Since(init_time)
		fmt.Printf("- Tempo de geração: %s\n", elapsed_time)

		// Mostrar informações sobre o módulo n
		fmt.Printf("- Módulo n gerado com %d bits\n", bbs.n.BitLen())

		// Gerar o número para demonstração
		fmt.Printf("- Gerando bits aleatórios...\n")
		randomNum := bbs.Next()

		// Exibir o tamanho real em bits
		bitLength := randomNum.BitLen()
		fmt.Printf("- Tamanho real: %d bits\n", bitLength)

		// Exibir a representação decimal
		fmt.Printf("- Valor decimal: %s\n", randomNum.String())

		// Exibir parcialmente a representação binária
		fmt.Printf("- Representação binária: %s\n", randomNum.Text(2))

		generatedNumbers[i] = randomNum
	}

	return bitSizes, generatedNumbers
}
