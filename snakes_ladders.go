package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/pkg/profile"
)

type gameBoard struct {
	featurePositions []int
	features         map[int]int
	knownPaths       map[int]int
	shortestPath     int
}

type RunFunc func(string)

var MAX_PATHS int = 1000000

func main() {
	typeP := os.Args[1]
	inputFile := os.Args[2]
	goodOrBad := os.Args[3]
	
	if typeP == "m" {
		defer profile.Start(profile.MemProfile).Stop()
	} else if typeP == "b" {
		defer profile.Start(profile.BlockProfile).Stop()
	} else if typeP == "t" {
		defer profile.Start(profile.TraceProfile).Stop()
	} else {
		defer profile.Start().Stop()
	}

	var rf RunFunc

	if goodOrBad == "g" {
		rf = RunFromFileGood
	} else {
		rf = RunFromFile
	}

	for i := 0; i < 10000; i++ {
		rf(inputFile)
	}
}

func RunFromFile(fileName string) {
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Println("Could not open file", err)
		return
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Scan()

	numBoards, err := strconv.Atoi(scanner.Text())
	if err != nil {
		fmt.Println("Failed to read number of boards")
		return
	}

	boards := make([]*gameBoard, numBoards)
	for i := 0; i < numBoards; i++ {
		boards[i] = &gameBoard{
			features:     make(map[int]int),
			knownPaths:   make(map[int]int),
			shortestPath: MAX_PATHS,
		}

		scanner.Scan()
		numLadders, err := strconv.Atoi(scanner.Text())
		if err != nil {
			fmt.Println("Could not read number of ladders")
			return
		}

		for x := 0; x < numLadders; x++ {
			scanner.Scan()
			coords := strings.Split(scanner.Text(), " ")

			ladder1, err := strconv.Atoi(coords[0])
			if err != nil {
				fmt.Println("Failed to read ladder coord 1")
			}

			ladder2, err := strconv.Atoi(coords[1])
			if err != nil {
				fmt.Println("Failed to read ladder coord 2")
			}

			boards[i].featurePositions = append(boards[i].featurePositions, ladder1)
			boards[i].features[ladder1] = ladder2
		}

		scanner.Scan()
		numSnakes, err := strconv.Atoi(scanner.Text())

		for x := 0; x < numSnakes; x++ {
			scanner.Scan()
			coords := strings.Split(scanner.Text(), " ")

			snake1, err := strconv.Atoi(coords[0])
			if err != nil {
				fmt.Println("Failed to read snake coord 1")
			}

			snake2, err := strconv.Atoi(coords[1])
			if err != nil {
				fmt.Println("Failed to read snake coord 2")
			}

			boards[i].featurePositions = append(boards[i].featurePositions, snake1)
			boards[i].features[snake1] = snake2
		}

		boards[i].featurePositions = append(boards[i].featurePositions, 1)
		boards[i].featurePositions = append(boards[i].featurePositions, 100)
		sort.Ints(boards[i].featurePositions)
	}

	resultChan := make(chan string, numBoards)
	boardChan := make(chan *gameBoard, numBoards)
	for b := 0; b < numBoards; b++ {
		boardChan <- boards[b]
	}

	for i := 0; i < numBoards; i++ {
		go func(index int) {
			board := <- boardChan
			board.findShortest(1, 0, board.featurePositions[1:], make(map[int]bool))

			result := ""
			if board.shortestPath < MAX_PATHS {
				result = fmt.Sprintf("Shortest path for board %d is %d", index, board.shortestPath)
			} else {
				result = fmt.Sprintf("No solution for board %d", index)
			}

			resultChan <- result
		}(i)
	}

	for x := 0; x < numBoards; x++ {
		fmt.Println(<- resultChan)
	}
}

func RunFromFileGood(fileName string) {
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Println("Could not open file", err)
		return
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Scan()

	numBoards, err := strconv.Atoi(scanner.Text())
	if err != nil {
		fmt.Println("Failed to read number of boards")
		return
	}

	boards := make([]*gameBoard, numBoards)
	for i := 0; i < numBoards; i++ {
		boards[i] = &gameBoard{
			features:     make(map[int]int),
			knownPaths:   make(map[int]int),
			shortestPath: MAX_PATHS,
		}

		scanner.Scan()
		numLadders, err := strconv.Atoi(scanner.Text())
		if err != nil {
			fmt.Println("Could not read number of ladders")
			return
		}

		for x := 0; x < numLadders; x++ {
			scanner.Scan()
			coords := strings.Split(scanner.Text(), " ")

			ladder1, err := strconv.Atoi(coords[0])
			if err != nil {
				fmt.Println("Failed to read ladder coord 1")
			}

			ladder2, err := strconv.Atoi(coords[1])
			if err != nil {
				fmt.Println("Failed to read ladder coord 2")
			}

			boards[i].featurePositions = append(boards[i].featurePositions, ladder1)
			boards[i].features[ladder1] = ladder2
		}

		scanner.Scan()
		numSnakes, err := strconv.Atoi(scanner.Text())

		for x := 0; x < numSnakes; x++ {
			scanner.Scan()
			coords := strings.Split(scanner.Text(), " ")

			snake1, err := strconv.Atoi(coords[0])
			if err != nil {
				fmt.Println("Failed to read snake coord 1")
			}

			snake2, err := strconv.Atoi(coords[1])
			if err != nil {
				fmt.Println("Failed to read snake coord 2")
			}

			boards[i].featurePositions = append(boards[i].featurePositions, snake1)
			boards[i].features[snake1] = snake2
		}

		boards[i].featurePositions = append(boards[i].featurePositions, 1)
		boards[i].featurePositions = append(boards[i].featurePositions, 100)
		sort.Ints(boards[i].featurePositions)
	}

	for i := 0; i < numBoards; i++ {
		board := boards[i]
		board.findShortest(1, 0, board.featurePositions[1:], make(map[int]bool))

		if board.shortestPath < MAX_PATHS {
			fmt.Printf("Shortest path for board %d is %d \n", i, board.shortestPath)
		} else {
			fmt.Printf("No solution for board %d \n", i)
		}
	}
}

func (b *gameBoard) findShortest(curPos int, runSum int, nextNodes []int, visited map[int]bool) {
	if _, ok := visited[curPos]; ok {
		return
	}

	if runSum >= b.shortestPath {
		return
	}

	newVisited := addToVisited(visited, curPos)

	// traverse the snake or ladder to get at new location
	var newNext []int
	if curPos > 1 {
		curPos = b.features[curPos]

		if curPos == 100 {
			if runSum < b.shortestPath {
				b.shortestPath = runSum
			}

			return
		}

		nextFeatIndex := sort.Search(len(b.featurePositions), func(i int) bool { return b.featurePositions[i] > curPos })
		newNext = b.featurePositions[nextFeatIndex:]
	} else {
		newNext = nextNodes
	}

	numNext := len(newNext)
	for i := 0; i < numNext; i++ {
		nextPos := newNext[i]

		if _, ok := visited[nextPos]; ok {
			continue
		}

		nextCost := b.findPathNum(curPos, nextPos)
		if nextCost < 0 {
			continue
		}

		nextSum := runSum + nextCost
		if nextPos != 100 {
			b.findShortest(nextPos, nextSum, newNext[i+1:], newVisited)
		} else if nextSum < b.shortestPath {
			b.shortestPath = nextSum
		}
	}
}

func (b *gameBoard) findPathNum(start int, end int) int {
	key := (100 * start) + end

	value, ok := b.knownPaths[key]

	if ok {
		return value
	}

	value = 0
	for current := start; current < end; {
		current += 6
		value++

		// Assumes there will never be six consecutive features
		if current != end {
			for i := 0; i < 6; i++ {
				if i == 5 {
					// back at original spot, no way past!
					value = -1
					b.knownPaths[key] = value
					return value
				}

				if _, ok := b.features[current]; !ok {
					break
				} else {
					current -= 1
				}
			}
		}
	}

	b.knownPaths[key] = value
	return value
}

// I hate this, must refactor
func addToVisited(orig map[int]bool, newValue int) map[int]bool {
	lenOrig := len(orig)

	newMap := make(map[int]bool, lenOrig+1)

	for k, v := range orig {
		newMap[k] = v
	}

	newMap[newValue] = true

	return newMap
}
