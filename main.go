package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

type sortable struct {
	index []int64
}

func (s sortable) Len() int {
	return len(s.index)
}

func (s sortable) Less(i, j int) bool {
	return s.index[i] < s.index[j]
}

func (s sortable) Swap(i, j int) {
	s.index[i], s.index[j] = s.index[j], s.index[i]
}

func main() {
	iterations := flag.Int("iterations", 3, "Max iterations")
	fileInput := flag.String("fileInput", "a_example.in", "Input filename")
	fileOutput := flag.String("fileOutput", "", "Output filename, appends to the input filename if empty")
	flag.Parse()

	target, numberss, err := ReadInput(*fileInput)
	if err != nil {
		log.Fatal(err)
	}

	rand.Seed(time.Now().Unix())
	bestSum := int64(0)
	bestPizzas := make([]int64, 0)

	for i := 0; i < *iterations; i++ {
		numbers := make([]int64, len(numberss))
		copy(numbers, numberss)

		sequence := make([]int64, len(numbers))
		for i := range sequence {
			sequence[i] = int64(i)
		}
		
		//Don't shuffle on 1 iteration
		if *iterations != 1 {
			rand.Shuffle(len(numbers), func(a, b int) {
				numbers[a], numbers[b] = numbers[b], numbers[a]
				sequence[a], sequence[b] = sequence[b], sequence[a]
			})
		}

		sum := int64(0)
		pizzas := make([]int64, 0)

		for j, n1 := range numbers {
			if n1+sum > target {
				break
			}

			sum += n1
			pizzas = append(pizzas, sequence[j])
		}

		if sum > bestSum {
			bestSum = sum
			bestPizzas = pizzas
		}

		if sum == target {
			fmt.Printf("Found precise in %v iterations\n", i)

			break
		}
	}

	if *fileOutput == "" {
		*fileOutput = *fileInput + ".out"
	}

	sort.Sort(sortable{bestPizzas})

	fmt.Printf("Target is %v, best is %v, delta is %v, the number of pizzas is %v\n", target, bestSum, target-bestSum, len(bestPizzas))

	err = SaveOutput(*fileOutput, target, bestSum, bestPizzas)
	if err != nil {
		log.Fatal(err)
	}
}

func ReadInput(filename string) (int64, []int64, error) {
	file, err := os.Open("data/" + filename)
	if err != nil {
		return 0, nil, err
	}

	b := make([]byte, 8*1024)
	input := make([]byte, 0)
	var size int

	for {
		size, err = file.Read(b)

		if err == io.EOF {
			break
		}

		input = append(input, b[:size]...)
	}

	input = append(input, b[:size]...)

	lines := bytes.Split(input, []byte("\n"))

	numbers := make([]int64, 0)

	var target int64

	for i, line := range lines {
		nums := bytes.Split(line, []byte(" "))

		for _, num := range nums {
			if string(num) == "" || string(num) == "\n" {
				break
			}

			parsed, err := strconv.ParseInt(string(num), 10, 64)
			if err != nil {
				return 0, nil, err
			}

			if i == 0 {
				target = parsed

				break
			}

			numbers = append(numbers, parsed)
		}
	}

	return target, numbers, nil
}

func SaveOutput(filename string, target int64, sum int64, pizzas []int64) error {
	file, err := os.OpenFile("data/"+filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0775)
	if err != nil {
		return err
	}

	pzz := make([]string, 0)

	for _, pizza := range pizzas {
		pzz = append(pzz, strconv.FormatInt(pizza, 10))
	}

	str := fmt.Sprintf("%v\n%s\n", len(pizzas), strings.Join(pzz, " "))

	_, err = file.Write([]byte(str))
	if err != nil {
		return err
	}

	return nil
}
