package main

import (
	"bufio"
	"fmt"
	"os"
	"sync"
)

func readFile(filename string) ([][]rune, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("wrong file format")
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	var worldSearch [][]rune
	for scanner.Scan() {
		line := scanner.Text()
		worldSearch = append(worldSearch, []rune(line))
	}
	return worldSearch, nil
}

func checkForXMAS(worldSearch [][]rune, i int, j int) int {
	xmasFound := 0
	lenX := len(worldSearch)
	lenY := len(worldSearch[0])
	if worldSearch[i][j] != 'X' {
		return 0
	}
	//checking verticals
	if i+3 < lenX && worldSearch[i+1][j] == 'M' && worldSearch[i+2][j] == 'A' && worldSearch[i+3][j] == 'S' {
		xmasFound += 1
	}
	if i-3 >= 0 && worldSearch[i-1][j] == 'M' && worldSearch[i-2][j] == 'A' && worldSearch[i-3][j] == 'S' {
		xmasFound += 1
	}
	//checking horizontals
	if j+3 < lenY && worldSearch[i][j+1] == 'M' && worldSearch[i][j+2] == 'A' && worldSearch[i][j+3] == 'S' {
		xmasFound += 1
	}
	if j-3 >= 0 && worldSearch[i][j-1] == 'M' && worldSearch[i][j-2] == 'A' && worldSearch[i][j-3] == 'S' {
		xmasFound += 1
	}
	//checking diagonals
	if i+3 < lenX && j+3 < lenY && worldSearch[i+1][j+1] == 'M' && worldSearch[i+2][j+2] == 'A' && worldSearch[i+3][j+3] == 'S' {
		xmasFound += 1
	}
	if i+3 < lenX && j-3 >= 0 && worldSearch[i+1][j-1] == 'M' && worldSearch[i+2][j-2] == 'A' && worldSearch[i+3][j-3] == 'S' {
		xmasFound += 1
	}
	if i-3 >= 0 && j+3 < lenY && worldSearch[i-1][j+1] == 'M' && worldSearch[i-2][j+2] == 'A' && worldSearch[i-3][j+3] == 'S' {
		xmasFound += 1
	}
	if i-3 >= 0 && j-3 >= 0 && worldSearch[i-1][j-1] == 'M' && worldSearch[i-2][j-2] == 'A' && worldSearch[i-3][j-3] == 'S' {
		xmasFound += 1
	}
	return xmasFound
}

func checkForCrossMAS(worldSearch [][]rune, i int, j int) int {
	crossMASFound := 0
	lenX := len(worldSearch)
	lenY := len(worldSearch[0])
	if worldSearch[i][j] != 'A' || i == 0 || j == 0 || i == lenX-1 || j == lenY-1 {
		return 0
	}
	if worldSearch[i-1][j-1] == 'M' && worldSearch[i+1][j+1] == 'S' && worldSearch[i+1][j-1] == 'M' && worldSearch[i-1][j+1] == 'S' {
		crossMASFound += 1
	}
	if worldSearch[i-1][j-1] == 'M' && worldSearch[i+1][j+1] == 'S' && worldSearch[i+1][j-1] == 'S' && worldSearch[i-1][j+1] == 'M' {
		crossMASFound += 1
	}
	if worldSearch[i-1][j-1] == 'S' && worldSearch[i+1][j+1] == 'M' && worldSearch[i+1][j-1] == 'M' && worldSearch[i-1][j+1] == 'S' {
		crossMASFound += 1
	}
	if worldSearch[i-1][j-1] == 'S' && worldSearch[i+1][j+1] == 'M' && worldSearch[i+1][j-1] == 'S' && worldSearch[i-1][j+1] == 'M' {
		crossMASFound += 1
	}	
	return crossMASFound
}

func part1(worldSearch [][]rune, result chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	xmasFound := 0
	lenX := len(worldSearch)
	lenY := len(worldSearch[0])
	for i := 0; i < lenX; i++ {
		for j := 0; j < lenY; j++ {
			xmasFound += checkForXMAS(worldSearch, i, j)
		}
	}
	result <- xmasFound
}

func part2(worldSearch [][]rune, result chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	crossMASFound := 0
	lenX := len(worldSearch)
	lenY := len(worldSearch[0])
	for i := 0; i < lenX; i++ {
		for j := 0; j < lenY; j++ {
			crossMASFound += checkForCrossMAS(worldSearch, i, j)
		}
	}
	result <- crossMASFound
}

func runFunc(part1Channel chan int, part2Channel chan int) {
	defer close(part1Channel)
	worldSearch, err := readFile("input.txt")
	if err != nil {
		return
	}
	wg := new(sync.WaitGroup)
	wg.Add(2)
	go part1(worldSearch, part1Channel, wg)
	go part2(worldSearch, part2Channel, wg)
	wg.Wait()
	return
}
func main() {
	xmasCounter := make(chan int)
	crossMasCounter := make(chan int)
	go runFunc(xmasCounter, crossMasCounter)
	fmt.Printf("Solution for Part 1:\nWe found %d apperances of XMAS\n\n", <-xmasCounter)
	fmt.Printf("Solution for Part 2:\nWe found %d apperances of XMAS\n\n", <-crossMasCounter)
	return
}
