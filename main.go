package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"
)

func main() {
	iterations := flag.Int("iterations", 1, "Max iterations")
	filename := flag.String("filename", "in1.txt", "Input filename")
	flag.Parse()

	target, numbers, err := ReadInput(*filename)
	if err != nil {
		log.Fatal(err)
	}

	rand.Seed(time.Now().Unix())
	bestSum := int64(0)
	bestSlices := make([]int64, 0)

	for i := 0; i < *iterations; i++ {
		//Don't shuffle on 1 iteration
		if *iterations != 1 {
			rand.Shuffle(len(numbers), func(a, b int) { numbers[a], numbers[b] = numbers[b], numbers[a] })
			//fmt.Println(numbers)
		}

		sum := int64(0)
		slices := make([]int64, 0)

		for _, n1 := range numbers {
			if n1+sum > target {
				break
			}

			sum += n1
			slices = append(slices, n1)
		}

		if sum > bestSum {
			bestSum = sum
			bestSlices = slices
		}

		if sum == target {
			fmt.Printf("Found precise in %v iterations\n", i)

			break
		}
	}

	fmt.Printf("Target is %v, best is %v, delta is %v, the number of pizzas is %v\n", target, bestSum, target-bestSum, len(bestSlices))
}

func ReadInput(filename string) (int64, []int64, error) {
	file, err := os.Open(filename)
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
