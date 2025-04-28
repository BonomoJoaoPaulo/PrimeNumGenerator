# PrimeNumGenerator

Esse repositório contém a implementação em Go de dois geradores 
 de números pseudo-aleatórios em Go:
- Blum Blum Shub
- Lagged Fibonacci Generator

E também dois testes de primalidade:
- Teste de Primalidade de Fermat
- Teste de Miller-Rabin

A implementação faz parte do Trabalho sobre Números Primos 
 da disciplina de Segurança em Computação (INE5429) do curso de
 Ciências da Computação da UFSC.

## Organização do Repositório
Os arquivos em _/prng_ referem-se às implementações dos geradores
 e os arquivos em _/pta_ às implementações dos testes de primalidade.

O script bash _run_tests.sh_ executa 10 vezes cada um dos dois geradores de números
 pseudo-aleatórios, então usa os valores gerados como entrada (cadidato) para os
 testes de primalidade, que então verificarão se aquele número é primo e, em caso
 negativo, gerará um número primo a partir do candidato fornecido.

Em _/output/_ estão duas subpastas, uma para cada um dos geradores de números, que
 trazem os arquivos de saída para 10 execuções do código, usados para o relatório do
 trabalho.
 
O arquivo _relatory.pdf_ traz um relatório com a descrição dos algoritmos
  escolhidos, os experimentos executados e uma análise do trabalho
  desenvolvido.

## Execução do código
### Pré-Requisitos
Este projeto foi desenvolvido utilizando apenas a biblioteca padrão do Go 
(não foram utilizadas bibliotecas externas).
Para compilar e executar o projeto, é necessário ter instalado a
 versão 1.18 ou superior do Go.

Você pode verificar a instalação do Go com:
```
go version
```
Caso precise instalar o Go, acesse: https://golang.org/dl/

### Executando
 Para executar o código usando o Blum Blum Shub como algoritmo gerador,
  e aplicar seus números gerados como candidatos para os dois testes de
  primalidade, você pode executar:
 ```
 go run main.go bbs
 ```
 Para o Lagged Fibonacci Generator:
 ```
 go run main.go fibonacci
 ```

 Alternativamente, caso queira rodar ambos 10 vezes, use o script pronto para Linux:
 ```
 ./run_test.sh
 ```

---
##### Última atualização em 28 de abril de 2025.
