// Esse arquivo traz a implementacao do algoritmo Lagged Fibonacci Generator
//  para gerar numeros pseudoaleatorios grandes.

package prng

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"time"
)

// LaggedFibonacciGenerator implementa o algoritmo de mesmo nome
//
//	para gerar os numeros pseudoaleatorios grandes.
type LaggedFibonacciGenerator struct {
	j, k     int
	state    []*big.Int
	size     int
	modValue *big.Int
	bitSize  int
}

// A funcao NewLFG cria um novo gerador com os parametros especificados
// size --> 	define o tamanho do buffer de estado
// j, k --> 	definem os indices usados na soma
// bitSize --> 	define o tamanho em bits dos numeros gerados
// returns --> 	retorna um ponteiro para o gerador
func NewLFG(size, j, k int, bitSize int) *LaggedFibonacciGenerator {
	// Garantimos que j < k
	if j >= k {
		panic("j deve ser menor que k")
	}

	if size < k {
		size = k // Garantimos que o buffer de estado seja pelo menos do tamanho de k
	}

	// Inicializamos o gerador
	lfg := &LaggedFibonacciGenerator{
		j:       j,
		k:       k,
		state:   make([]*big.Int, size),
		size:    size,
		bitSize: bitSize,
	}

	// Definimos o modValue como 2^bitSize
	lfg.modValue = new(big.Int).Lsh(big.NewInt(1), uint(bitSize))

	// Inicializamos o estado com valores aleatorios verdadeiros do tamanho apropriado
	for i := 0; i < size; i++ {
		// Criamos um numero aleatorio criptograficamente seguro com o tamanho de bits desejado
		randBits, err := rand.Int(rand.Reader, new(big.Int).Sub(lfg.modValue, big.NewInt(1)))
		if err != nil {
			// Fallback para um metodo menos seguro se rand.Int falhar
			randBits = generateFallbackRandom(bitSize)
		}

		// Aqui garantimos que o numero tem um tamanho proximo ao desejado
		// Definimos o bit mais significativo para garantir o tamanho minimo
		if bitSize > 1 {
			randBits.SetBit(randBits, bitSize-1, 1)
			// Configuramos tambem alguns bits aleatorios para garantir certa variacao
			randBits.SetBit(randBits, bitSize/2, uint(time.Now().UnixNano()%2))
			randBits.SetBit(randBits, bitSize/3, uint((time.Now().UnixNano()+int64(i))%2))
		}

		lfg.state[i] = randBits
	}

	return lfg
}

// A funcao generateFallbackRandom gera um numero aleatorio grande usando um
// metodo menos seguro mas mais garantido de funcionar em todos os ambientes
// caso o rand.Int usado em NewLFG falhe.
func generateFallbackRandom(bitSize int) *big.Int {
	result := new(big.Int)

	// Calculamos quantos blocos de 31 bits precisamos
	// (usamos 31 em vez de 32 para evitar problemas com sinal em int64 no Go)
	blocks := (bitSize + 30) / 31

	// Geracao de cada bloco
	for i := 0; i < blocks; i++ {
		// Geracao de um numero aleatorio de 31 bits usando o timestamp
		seed := time.Now().UnixNano() + int64(i*9999)
		blockValue := new(big.Int).Mod(big.NewInt(seed), big.NewInt(1<<31))

		// Deslocamos o bloco para sua posição correta
		if i > 0 {
			blockValue.Lsh(blockValue, uint(i*31))
		}

		// Adicionamos ao resultado
		result.Or(result, blockValue)

		time.Sleep(time.Nanosecond)
	}

	return result
}

// Next gera e retorna o proximo numero na sequencia pseudoaleatoria
//
//	e atualiza o estado do gerador. O resultado eh um ponteiro
//	para um big.Int que representa o numero gerado.
//
// O numero gerado eh o resultado da soma dos dois numeros anteriores
//
//	na sequencia, com os indices j e k definidos no construtor.
func (lfg *LaggedFibonacciGenerator) Next() *big.Int {
	// Calculamos o proximo valor como state[i-j] + state[i-k] mod 2^bitSize
	result := new(big.Int)

	// state[i-j] + state[i-k]
	result.Add(lfg.state[lfg.size-lfg.j], lfg.state[lfg.size-lfg.k])
	result.Mod(result, lfg.modValue)

	// Deslocamos todos os valores no array
	for i := 0; i < lfg.size-1; i++ {
		lfg.state[i] = lfg.state[i+1]
	}

	// Adicionamos o novo valor ao final
	lfg.state[lfg.size-1] = result

	return new(big.Int).Set(result)
}

func Lfg() ([]int, []*big.Int) {
	// Tamanhos de bits para testar especificados no enunciado do trabalho
	bitSizes := []int{40, 56, 80, 128, 168, 224, 256, 512, 1024, 2048, 4096}
	generatedNumbers := make([]*big.Int, len(bitSizes))

	fmt.Println("Gerando números pseudoaleatórios com Lagged Fibonacci Generator")
	fmt.Println("=============================================================")

	// Usamos j=7, k=10 como exemplo de parametros comuns para LFG
	// usando como ref. o segundo volume da serie de livros
	// The Art of Computer Programming
	j := 7
	k := 10

	for i, bits := range bitSizes {
		fmt.Printf("\nGerando número de %d bits:\n", bits)

		// Criamos um novo gerador para cada tamanho de bits
		lfg := NewLFG(k, j, k, bits)

		// "Aquecemos" o gerador descartando alguns valores iniciais
		for i := 0; i < 20; i++ {
			lfg.Next()
		}

		startTime := time.Now()

		// Geramos o numero
		randomNum := lfg.Next()

		elapsedTime := time.Since(startTime)
		fmt.Printf("- Tempo de geração: %s\n", elapsedTime)

		// Exibimos o tamanho real em bits do numero gerado
		bitLength := randomNum.BitLen()
		fmt.Printf("- Tamanho real: %d bits\n", bitLength)

		// Exibimos a representacao decimal
		fmt.Printf("- Valor decimal: %s\n", randomNum.String())

		// Exibimos a representacao binaria
		fmt.Printf("- Representação binária: %s\n", randomNum.Text(2))

		// Por fim, verificamos se o numero tem o tamanho esperado ou proximo disso (dentro de 4 bits)
		if bitLength < bits-4 {
			fmt.Printf("AVISO: O número gerado tem menos bits que o solicitado (%d < %d)\n", bitLength, bits)
		}

		generatedNumbers[i] = randomNum
	}

	return bitSizes, generatedNumbers
}
