package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
)

func readFile(filename string) ([]int, []int, map[int]int, error) {
	var leftList, rightList []int
	freq := make(map[int]int)
	file, err := os.Open(filename)
	if err != nil {
		return nil, nil, nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var lineNumber int
	lineNumber = 0
	for scanner.Scan() {
		lineNumber += 1
		line := scanner.Text()
		numbers := strings.Split(line, "   ")
		if len(numbers) != 2 {
			return nil, nil, nil, fmt.Errorf("wrong file format in line: %d", lineNumber)
		}

		number1, err := strconv.Atoi(numbers[0])
		if err != nil {
			return nil, nil, nil, fmt.Errorf("wrong number format in line %d", lineNumber)
		}
		number2, err := strconv.Atoi(numbers[1])
		if err != nil {
			return nil, nil, nil, fmt.Errorf("wrong number format in line %d", lineNumber)
		}
		leftList = append(leftList, number1)
		rightList = append(rightList, number2)
		freq[number2] += 1
	}
	return leftList, rightList, freq, nil
}

// func readFileLists(filename string) ([]int, []int, error) {
// 	var list1, list2 []int
// 	file, err := os.Open(filename)
// 	if err != nil {
// 		return nil, nil, err
// 	}
// 	defer file.Close()

// 	scanner := bufio.NewScanner(file)
// 	var lineNumber int
// 	lineNumber = 0
// 	for scanner.Scan() {
// 		lineNumber += 1
// 		line := scanner.Text()
// 		numbers := strings.Split(line, "   ")
// 		// fmt.Println(numbers)
// 		// fmt.Println(len(numbers))
// 		if len(numbers) != 2 {
// 			return nil, nil, fmt.Errorf("wrong file format in line: %d", lineNumber)
// 		}

// 		number1, err := strconv.Atoi(numbers[0])
// 		if err != nil {
// 			return nil, nil, fmt.Errorf("wrong number format in line %d", lineNumber)
// 		}
// 		number2, err := strconv.Atoi(numbers[1])
// 		if err != nil {
// 			return nil, nil, fmt.Errorf("wrong number format in line %d", lineNumber)
// 		}

// 		list1 = append(list1, number1)
// 		list2 = append(list2, number2)
// 	}

// 	return list1, list2, nil
// }

// func readFileListMap(filename string) ([]int, map[int]int, error) {
// 	var list []int
// 	freq := make(map[int]int)
// 	file, err := os.Open(filename)
// 	if err != nil {
// 		return nil, nil, err
// 	}
// 	defer file.Close()

// 	scanner := bufio.NewScanner(file)
// 	var lineNumber int
// 	lineNumber = 0
// 	for scanner.Scan() {
// 		lineNumber += 1
// 		line := scanner.Text()
// 		numbers := strings.Split(line, "   ")
// 		// fmt.Println(numbers)
// 		// fmt.Println(len(numbers))
// 		if len(numbers) != 2 {
// 			return nil, nil, fmt.Errorf("wrong file format in line: %d", lineNumber)
// 		}

// 		number1, err := strconv.Atoi(numbers[0])
// 		if err != nil {
// 			return nil, nil, fmt.Errorf("wrong number format in line %d", lineNumber)
// 		}
// 		number2, err := strconv.Atoi(numbers[1])
// 		if err != nil {
// 			return nil, nil, fmt.Errorf("wrong number format in line %d", lineNumber)
// 		}

// 		list = append(list, number1)

// 		freq[number2] += 1
// 	}

// 	return list, freq, nil
// }

func absDiff(number1, number2 int) int {
	result := number1 - number2
	if result < 0 {
		return result * -1
	}
	return result
}

func part1(leftList []int, rightList []int, result chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	var totalDifference int
	totalDifference = 0
	for i := 0; i < len(leftList); i++ {
		difference := absDiff(leftList[i], rightList[i])
		totalDifference += difference
	}
	result <- totalDifference
	return
}

func part2(leftList []int, freq map[int]int, result chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	var similarityScore int
	similarityScore = 0
	for i := 0; i < len(leftList); i++ {
		similarityScore += leftList[i] * freq[leftList[i]]
	}
	result <- similarityScore
	return
}

func runFunc(part1Channel chan int, part2Channel chan int) {
	defer close(part1Channel)
	defer close(part2Channel)
	leftList, rightList, freq, err := readFile("input.txt")
	if err != nil {
		return
	}
	sort.Ints(leftList)
	sort.Ints(rightList)
	wg := new(sync.WaitGroup)
	wg.Add(2)
	go part1(leftList, rightList, part1Channel, wg)
	go part2(leftList, freq, part2Channel, wg)
	wg.Wait()
	return
}

func main() {
	totalDifferenceResult := make(chan int)
	// go part1(totalDifferenceResult)
	similarityScoreResult := make(chan int)
	// go part2(similarityScoreResult)
	go runFunc(totalDifferenceResult, similarityScoreResult)
	fmt.Printf("Solution for Part 1:\nThe Total Difference is %d\n\n", <-totalDifferenceResult)
	fmt.Printf("Solution for Part 2:\nThe Similarity Score is %d\n", <-similarityScoreResult)
	return
}
