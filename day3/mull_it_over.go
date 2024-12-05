package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"regexp"
	"sync"
)

func readFile(filename string) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", fmt.Errorf("wrong file format")
	}
	defer file.Close()
	scanner :=	bufio.NewScanner(file)
	scanner.Scan()
	text := scanner.Text()
	return text, nil
}

func part2(muls []string, result chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	
	res := 0
	enabled := true
	for i, mul := range(muls) {
		isDoDont := strings.Index(mul, "d")
		if isDoDont >= 0 {
			isDont := strings.Index(mul, "'")
			if isDont >= 0 {
				enabled = false
				continue
			} else {
				enabled = true
				continue
			}
		}
		if !enabled {
			continue
		}
		commaIndex := strings.Index(mul, ",")
		closedBracketIndex := strings.Index(mul, ")")
		number1, err := strconv.Atoi(mul[4:commaIndex])
		if err != nil {
			fmt.Printf("something went wrong while finding the first opperand of the %d mul in part2\n", i+1)
			return
		}
		number2, err := strconv.Atoi(mul[commaIndex+1:closedBracketIndex])
		if err != nil {
			fmt.Printf("something went wrong while finding the second opperand of the %d mul in part2\n", i+1)
			return
		}
		res += number1 * number2
	}
	result <- res
	return 
}
func part1(muls []string, result chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	
	res := 0
	for i, mul := range(muls) {
		isDo := strings.Index(mul, "d")
		if isDo >= 0 {
			continue
		}
		commaIndex := strings.Index(mul, ",")
		closedBracketIndex := strings.Index(mul, ")")
		number1, err := strconv.Atoi(mul[4:commaIndex])
		if err != nil {
			fmt.Printf("something went wrong while finding the first opperand of the %d mul in part1\n", i+1)
			return
		}
		number2, err := strconv.Atoi(mul[commaIndex+1:closedBracketIndex])
		if err != nil {
			fmt.Printf("something went wrong while finding the second opperand of the %d mul in part1\n", i+1)
			return
		}
		res += number1 * number2
	}
	result <- res
	return 
}

func runFunc(part1Channel chan int, part2Channel chan int) {
	defer close(part1Channel)
	memory, err := readFile("input.txt")
	if err != nil {
		return
	}
	pattern := `mul\(\d{1,3},\d{1,3}\)|do\(\)|don't\(\)`
	regex:= regexp.MustCompile(pattern)
	muls := regex.FindAllString(memory, -1)
	wg := new(sync.WaitGroup)
	wg.Add(2)
	go part1(muls, part1Channel, wg)
	go part2(muls, part2Channel, wg)
	wg.Wait()
	return
}

func main() {
	part1Result := make(chan int)
	part2Result := make(chan int)
	go runFunc(part1Result, part2Result)
	fmt.Printf("Solution for Part 1:\nAdding up the results of all the multiplications we get %d\n\n", <-part1Result)
	fmt.Printf("Solution for Part 2:\nAdding up the results of all the enabled multiplications we get %d\n\n", <-part2Result)
	// fmt.Printf("Solution for Part 2\nThere are %d safe reports after dampening\n\n", <-safeDampenedReports)
}