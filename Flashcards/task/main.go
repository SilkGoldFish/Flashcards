package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type Card struct {
	Term       string
	Definition string
	Mistakes   int
}

func readLine() string {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	return scanner.Text()
}

func readAndLog() string {
	line := readLine()
	logBuffer = append(logBuffer, line)
	return line
}

func printAndLog(line string) {
	fmt.Println(line)
	logBuffer = append(logBuffer, line)
}

func add(FlashCards map[string]*Card) {
	printAndLog("The card")
	term := readAndLog()
	for _, ok := FlashCards[term]; ok; _, ok = FlashCards[term] {
		printAndLog("This term already exists. Try again:")
		term = readAndLog()
	}
	printAndLog("The definition of the card")
	definition := readAndLog()
	for ok := true; ok; {
		ok = false
		for _, v := range FlashCards {
			if v.Definition == definition {
				printAndLog("This definition already exists. Try again:")
				definition = readAndLog()
				ok = true
				break
			}
		}
	}
	FlashCards[term] = &Card{term, definition, 0}
	printAndLog(fmt.Sprintf("The pair (\"%s\":\"%s\") has been added.", term, definition))
}

func remove(FlashCards map[string]*Card) {
	printAndLog("Which card?")
	term := readAndLog()
	if _, ok := FlashCards[term]; ok {
		delete(FlashCards, term)
		printAndLog("The card has been removed.")
	} else {
		printAndLog(fmt.Sprintf("Can't remove \"%s\": there is no such card.", term))
	}
}

func readFile(fileName string, FlashCards map[string]*Card) {
	data, err := os.ReadFile(fileName)
	if len(data) == 0 {
		printAndLog("File not found.")
		return
	}
	if err != nil {
		log.Fatal(err)
	}
	var flashCards map[string]*Card
	err = json.Unmarshal([]byte(data), &flashCards)
	if err != nil {
		log.Fatal(err)
	}
	for k, v := range flashCards {
		FlashCards[k] = v
	}
	printAndLog(fmt.Sprintf("%d cards have been loaded.", len(flashCards)))
}

func importFile(FlashCards map[string]*Card) {
	printAndLog("File name:")
	fileName := readAndLog()
	readFile(fileName, FlashCards)
}

func importFrom(importFileName string, FlashCards map[string]*Card) {
	if importFileName == "" {
		return
	}
	readFile(importFileName, FlashCards)
}

func writeFile(fileName string, FlashCards map[string]*Card) {
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_RDWR, 0744)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	data, err := os.ReadFile(fileName)
	if err != nil {
		log.Fatal(err)
	}
	if len(data) != 0 {
		var flashCards map[string]*Card
		err = json.Unmarshal([]byte(data), &flashCards)
		if err != nil {
			log.Fatal(err)
		}
		for k, v := range flashCards {
			if _, ok := FlashCards[k]; !ok {
				FlashCards[k] = v
			}
		}
	}
	FlashCardsJSON, err := json.Marshal(FlashCards)
	if err != nil {
		log.Fatal(err)
	}
	_, err = fmt.Fprintf(file, "%s", string(FlashCardsJSON))
	if err != nil {
		log.Fatal(err)
	}
	printAndLog(fmt.Sprintf("%d cards have been saved.", len(FlashCards)))

}

func exportTo(exportFileName string, FlashCards map[string]*Card) {
	if exportFileName == "" {
		return
	}
	writeFile(exportFileName, FlashCards)
}

func exportFile(FlashCards map[string]*Card) {
	printAndLog("File name:")
	fileName := readAndLog()
	writeFile(fileName, FlashCards)
}

func ask(FlashCards map[string]*Card) {
	printAndLog("How many times to ask?")
	times, _ := strconv.Atoi(readAndLog())
	for {
		for k, v := range FlashCards {
			printAndLog(fmt.Sprintf("Print the definition of \"%s\":", k))
			anwser := readAndLog()
			if anwser == v.Definition {
				printAndLog("Correct!")
			} else {
				FlashCards[k].Mistakes++
				ok := true
				for key, val := range FlashCards {
					if anwser == val.Definition {
						ok = false
						printAndLog(fmt.Sprintf("Wrong. The right answer is \"%s\", but your definition is correct for \"%s\".", v.Definition, key))
						break
					}
				}
				if ok {
					printAndLog(fmt.Sprintf("Wrong. The right answer is \"%s\".", v.Definition))
				}
			}
			times--
			if times == 0 {
				return
			}
		}
	}
}

func logToFile() {
	printAndLog("File name:")
	fileName := readAndLog()
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0744)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	for _, line := range logBuffer {
		_, err = fmt.Fprintln(file, line)
		if err != nil {
			log.Fatal(err)
		}
	}
	printAndLog("The log has been saved.")
}

func hardestCard(FlashCards map[string]*Card) {
	maxMistakeNum := 0
	hardestCards := make([]string, 0)
	for _, v := range FlashCards {
		if v.Mistakes == 0 {
			continue
		}
		if v.Mistakes > maxMistakeNum {
			maxMistakeNum = v.Mistakes
			hardestCards = hardestCards[:0]
			hardestCards = append(hardestCards, v.Term)
		} else if v.Mistakes == maxMistakeNum {
			hardestCards = append(hardestCards, v.Term)
		}
	}
	if len(hardestCards) == 0 {
		printAndLog("There are no cards with errors.")
	} else if len(hardestCards) == 1 {
		printAndLog(fmt.Sprintf("The hardest card is \"%s\". You have %d errors answering it.", hardestCards[0], maxMistakeNum))
	} else {
		var printInfo strings.Builder
		for i, term := range hardestCards {
			printInfo.WriteString("\"" + term + "\"")
			if i != len(hardestCards)-1 {
				printInfo.WriteString(", ")
			} else {
				printInfo.WriteString(". You have " + strconv.Itoa(maxMistakeNum) + " errors answering it.")
			}
		}
		printAndLog("The hardest cards are " + printInfo.String())
	}
}

func resetStats(FlashCards map[string]*Card) {
	for _, v := range FlashCards {
		v.Mistakes = 0
	}
	printAndLog("Card statistics have been reset.")
}

var logBuffer = make([]string, 0)

func main() {
	FlashCards := make(map[string]*Card)
	importFileName := flag.String("import_from", "", "Enter import file name")
	exportFileName := flag.String("export_to", "", "Enter export file name")
	flag.Parse()
	importFrom(*importFileName, FlashCards)
	for action := ""; action != "exit"; {
		printAndLog("Input the action (add, remove, import, export, ask, exit, log, hardest card, reset stats):")
		action = readAndLog()
		switch action {
		case "add":
			add(FlashCards)
		case "remove":
			remove(FlashCards)
		case "import":
			importFile(FlashCards)
		case "export":
			exportFile(FlashCards)
		case "ask":
			ask(FlashCards)
		case "log":
			logToFile()
		case "hardest card":
			hardestCard(FlashCards)
		case "reset stats":
			resetStats(FlashCards)
		case "exit":
			printAndLog("Bye bye!")
			break
		}
	}
	exportTo(*exportFileName, FlashCards)
}
