package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

func main() {

	numbers := GenerateNumbers(100000000)

	tamanho := len(numbers)
	fmt.Printf("Tamanho do array: %d \n", tamanho)

	t := time.Now()
	soma := Add(numbers)
	fmt.Printf("Soma dos números gerados: %d em %s \n", soma, time.Since(t))

	t = time.Now()
	soma = AddConcurrent(numbers)
	fmt.Printf("Soma dos números gerados: %d em %s \n", soma, time.Since(t))

}

func Add(numbers []int) int64 {
	var sum int64
	for _, n := range numbers {
		sum += int64(n)
	}
	return sum
}

func AddConcurrent(numbers []int) int64 {

	// Recuperar q quantidade de processadores e configurar o Go para Utilizar todos
	numOfCores := runtime.NumCPU()
	runtime.GOMAXPROCS(numOfCores)

	// Variável que armazenará a soma
	var sum int64

	// Tamanho do array (quantidade de números)
	max := len(numbers)
	// fmt.Printf("Tamanho do array: %d \n", max)

	// sizeOfPart é o tamanho do array dividido entre os processadores
	sizeOfPart := max / numOfCores
	// fmt.Printf("Tamanho do array dividido entre os processadores: %d \n", sizeOfPart)

	// WaitGroup para aguardar o término de todos os processadores
	var wg sync.WaitGroup

	// Para cada processador, será iniciado um novo goroutine que processará uma parte do array principal
	for i := 0; i < numOfCores; i++ {

		// fmt.Printf("Processador %d iniciado \n", i)

		// Dividir a entrada em partes para processar por cada processador
		start := i * sizeOfPart
		// fmt.Printf("Início do processador %d: %d \n", i, start)
		end := start + sizeOfPart
		// fmt.Printf("Fim do processador %d: %d \n", i, end)
		part := numbers[start:end]

		wg.Add(1)

		go func(nums []int) {
			defer wg.Done()

			var partSum int64

			// Somar os números da parte
			for _, n := range nums {
				partSum += int64(n)
			}

			// Adicionar a soma da parte à soma total
			atomic.AddInt64(&sum, partSum)

		}(part)

	}
	wg.Wait()
	return sum

}

func GenerateNumbers(qtde int) []int {
	numbers := make([]int, 0, qtde) // Cria um array de inteiros com o tamanho definido (importante usar o make(tipo, 0, tamanho))!
	for i := int(0); i < qtde; i++ {
		aleatorio := rand.Intn(qtde)
		numbers = append(numbers, aleatorio)
	}
	return numbers
}
