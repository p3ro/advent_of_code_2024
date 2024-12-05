package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// TODO: refactor code so we make better use of the channels and we don't read the file and create the same slices two times
func safeReportsWithoudDampening(filename string) (int, error) {
	file, err := os.Open(filename)
	if err != nil {
		return -1, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var lineNumber, safeReports int
	lineNumber = 0
	safeReports = 0
	for scanner.Scan() {
		lineNumber += 1
		line := scanner.Text()
		levels := strings.Split(line, " ")
		// fmt.Println(levels)
		// fmt.Println(len(levels))
		badLevel, err := checkLevels(levels)
		if err != nil {
			fmt.Printf("wrong number format in line %d\n", lineNumber)
			return -1, err
		}
		if badLevel == -1 {
			safeReports += 1
		}
	}

	return safeReports, nil
}

func safeReportsWithDampening(filename string) (int, error) {
	file, err := os.Open(filename)
	if err != nil {
		return -1, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var lineNumber, safeReports int
	lineNumber = 0
	safeReports = 0
	for scanner.Scan() {
		lineNumber += 1
		line := scanner.Text()
		levels := strings.Split(line, " ")
		// fmt.Println(levels)
		// fmt.Println(len(levels))
		badLevel, err := checkLevels(levels)
		if err != nil {
			fmt.Printf("wrong number format in line %d\n", lineNumber)
			return -1, err
		}
		if badLevel == -1 {
			safeReports += 1
			continue
		}
		var check int
		if badLevel != 0 { //checking the first level anyway cause otherwise it would dodge our test in something like 48 46 47 49 51 54 56
			//where we would only consider removing 46 47 48 otherwise
			levelsWithoutFirstLevel := make([]string, 0, len(levels)-1)
			levelsWithoutFirstLevel = append(levelsWithoutFirstLevel, levels[1:]...)
			check, err = checkLevels(levelsWithoutFirstLevel)
			if err != nil {
				fmt.Printf("wrong number format in line %d\n", lineNumber)
				return -1, err
			}
			if check == -1 {
				safeReports += 1
				continue
			}
		}
		levelsWithoutBadLevel1 := make([]string, 0, len(levels)-1)
		levelsWithoutBadLevel1 = append(levelsWithoutBadLevel1, levels[:badLevel]...)
		levelsWithoutBadLevel1 = append(levelsWithoutBadLevel1, levels[badLevel+1:]...)
		check, err = checkLevels(levelsWithoutBadLevel1)
		if err != nil {
			fmt.Printf("wrong number format in line %d\n", lineNumber)
			return -1, err
		}
		if check == -1 {
			safeReports += 1
			continue
		}
		levelsWithoutBadLevel2 := make([]string, 0, len(levels)-1)
		levelsWithoutBadLevel2 = append(levelsWithoutBadLevel2, levels[:badLevel+1]...)
		levelsWithoutBadLevel2 = append(levelsWithoutBadLevel2, levels[badLevel+2:]...)
		check, err = checkLevels(levelsWithoutBadLevel2)
		if err != nil {
			fmt.Printf("wrong number format in line %d\n", lineNumber)
			return -1, err
		}
		if check == -1 {
			safeReports += 1
		}
		if badLevel == len(levels)-2 {
			continue
		}
		levelsWithoutBadLevel3 := make([]string, 0, len(levels)-1)
		levelsWithoutBadLevel3 = append(levelsWithoutBadLevel3, levels[:badLevel+2]...)
		levelsWithoutBadLevel3 = append(levelsWithoutBadLevel3, levels[badLevel+3:]...)
		check, err = checkLevels(levelsWithoutBadLevel3)
		if err != nil {
			fmt.Printf("wrong number format in line %d\n", lineNumber)
			return -1, err
		}
		if check == -1 {
			safeReports += 1
		}
	}

	return safeReports, nil
}

func checkLevels(levels []string) (int, error) {
	isIncreasing := 0 // if isIncreasing is 0 the value is not set, if it's -1 it's false and if it's 1 it's true
	isSafe := true
	prevLevel, err := strconv.Atoi(levels[0])
	if err != nil {
		return -1, err
	}
	badLevel := -1
	for i := 1; i < len(levels) && isSafe; i++ {
		currLevel, err := strconv.Atoi(levels[i])
		if err != nil {
			return -1, err
		}
		if !isGradual(currLevel, prevLevel) {
			isSafe = false
			badLevel = i - 1
		}

		//if it's the first two levels we check to see if the values are increasing or decreasing
		if isIncreasing == 0 {
			isIncreasing = boolToInt(prevLevel < currLevel)
		}

		//for the next levels we check to see if the values follow the same direction as the previous
		if isIncreasing*boolToInt(prevLevel < currLevel) == -1 {
			isSafe = false
			badLevel = i - 1
		}
		prevLevel = currLevel
	}

	return badLevel, nil
}

func isGradual(currLevel int, prevLevel int) bool {
	diff := absDiff(currLevel, prevLevel)
	return (diff >= 1 && diff <= 3)
}

func boolToInt(b bool) int {
	var i int
	if b {
		i = 1
	} else {
		i = -1
	}
	return i
}

func absDiff(number1, number2 int) int {
	result := number1 - number2
	if result < 0 {
		return result * -1
	}
	return result
}

func part1(result chan int) {
	defer close(result)
	safeReports, err := safeReportsWithoudDampening("input.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	result <- safeReports
	return
}

func part2(result chan int) {
	defer close(result)
	safeReports, err := safeReportsWithDampening("input.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	result <- safeReports
	return
}

func main() {
	safeReports := make(chan int)
	go part1(safeReports)
	safeDampenedReports := make(chan int)
	go part2(safeDampenedReports)
	fmt.Printf("Solution for Part 1:\nThere are %d safe reports\n\n", <-safeReports)
	fmt.Printf("Solution for Part 2\nThere are %d safe reports after dampening\n\n", <-safeDampenedReports)
}
