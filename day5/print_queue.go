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

type Pair struct {
	page1, page2 int
}

func readFile(filename string) (map[Pair]bool, [][]int, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, nil, fmt.Errorf("wrong file format")
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	orderingMap := make(map[Pair]bool)
	var queries [][]int
	readingOrdering := true
	lineNumber := 0
	for scanner.Scan() {
		lineNumber += 1
		line := scanner.Text()
		if line == "" {
			readingOrdering = false
			continue
		}
		if readingOrdering {
			numbers := strings.Split(line, "|")
			if len(numbers) != 2 {
				return nil, nil, fmt.Errorf("wrong file format in line: %d", lineNumber)
			}

			page1, err := strconv.Atoi(numbers[0])
			if err != nil {
				return nil, nil, fmt.Errorf("wrong number format in line %d", lineNumber)
			}
			page2, err := strconv.Atoi(numbers[1])
			if err != nil {
				return nil, nil, fmt.Errorf("wrong number format in line %d", lineNumber)
			}
			truePair := Pair{page1, page2}
			falsePair := Pair{page2, page1}
			orderingMap[truePair] = true
			orderingMap[falsePair] = false
		} else {
			numbers := strings.Split(line, ",")
			pages := make([]int, 0, len(numbers))
			for _, number := range numbers {
				page, err := strconv.Atoi(number)
				if err != nil {
					return nil, nil, fmt.Errorf("wrong file format in line: %d", lineNumber)
				}
				pages = append(pages, page)
			}
			queries = append(queries, pages)
		}
	}
	return orderingMap, queries, nil
}

func orderedPages(orderingMap map[Pair]bool, pages []int) bool {
	for i := 0; i < len(pages); i++ {
		for j := i + 1; j < len(pages); j++ {
			checkPair := Pair{pages[i], pages[j]}
			if isOrdered, ok := orderingMap[checkPair]; ok && !isOrdered {
				return false
			}
		}
	}
	return true
}

func reorderPages(orderingMap map[Pair]bool, pages []int) []int {
	newPages := make([]int, 0, len(pages))
	newPages = append(newPages, pages...)
	sort.Slice(newPages, func(i, j int) bool {
		return orderingMap[Pair{newPages[i], newPages[j]}]
	})
	return newPages
}

func reorderedPages(orderingMap map[Pair]bool, pages []int) (bool, []int) {
	for i := 0; i < len(pages); i++ {
		for j := i + 1; j < len(pages); j++ {
			checkPair := Pair{pages[i], pages[j]}
			if isOrdered, ok := orderingMap[checkPair]; ok && !isOrdered {
				return true, reorderPages(orderingMap, pages)
			}
		}
	}
	return false, nil
}

func part1(orderingMap map[Pair]bool, queries [][]int, result chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	sumOfMiddles := 0
	for _, pages := range queries {
		if orderedPages(orderingMap, pages) {
			middleIndex := len(pages) / 2
			sumOfMiddles += pages[middleIndex]
		}
	}
	result <- sumOfMiddles
	return
}

func part2(orderingMap map[Pair]bool, queries [][]int, result chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	sumOfMiddles := 0
	for _, pages := range queries {
		if reorder, newPages := reorderedPages(orderingMap, pages); reorder {
			middleIndex := len(newPages) / 2
			sumOfMiddles += newPages[middleIndex]
		}
	}
	result <- sumOfMiddles
}

func runFunc(part1Channel chan int, part2Channel chan int) {
	defer close(part1Channel)
	defer close(part2Channel)
	orderingMap, queries, err := readFile("input.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	// fmt.Println(orderingMap)
	// fmt.Printf("\n\n\n\n")
	// fmt.Println(queries)
	wg := new(sync.WaitGroup)
	wg.Add(2)
	go part1(orderingMap, queries, part1Channel, wg)
	go part2(orderingMap, queries, part2Channel, wg)
	wg.Wait()
	part2Channel <- 2
	return
}

func main() {
	correctQueries := make(chan int)
	reorderedQueries := make(chan int)
	go runFunc(correctQueries, reorderedQueries)
	fmt.Printf("Solution for Part 1:\nThe sum of the middle numbers of the correct queries is %d\n\n", <-correctQueries)
	fmt.Printf("Solution for Part 2:\nThe sum of the middle numbers of the reordered queries is %d\n\n", <-reorderedQueries)
	return
}
