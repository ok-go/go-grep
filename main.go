package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
)

var (
	pattern       string
	filename      string
	afterContext  int
	beforeContext int

	found = 1
)

func init() {
	var context int
	flag.IntVar(&context, "C", 0, "Print num lines of leading and trailing context surrounding each match.")
	flag.IntVar(&afterContext, "A", 0, "Print num lines of trailing context after each match.")
	flag.IntVar(&beforeContext, "B", 0, "Print num lines of leading context before each match.")
	flag.Parse()

	if context > 0 {
		afterContext = context
		beforeContext = context
	}

	pattern = flag.Arg(0)
	filename = flag.Arg(1)
}

func read(path string) {
	file, err := os.Open(path)
	if err != nil {
		println(err)
		os.Exit(-1)
	}
	defer file.Close()

	buff := make([]string, 0, beforeContext)
	scanner := bufio.NewScanner(file)
	for lineNumber := 0; scanner.Scan(); lineNumber++ {
		if strings.Contains(scanner.Text(), pattern) {
			if len(buff) > 0 && lineNumber > 0 {
				fmt.Println("...")
			}
			for _, line := range buff {
				printLine(line)
			}
			printLine(scanner.Text())
			for i := 0; i < afterContext; i++ {
				if scanner.Scan() {
					lineNumber++
					printLine(scanner.Text())
					if strings.Contains(scanner.Text(), pattern) {
						i = -1
					}
				}
			}

			found = 0
			continue
		}

		if len(buff) < beforeContext {
			buff = append(buff, scanner.Text())
		} else if beforeContext > 0 {
			buff = append(buff[1:], scanner.Text())
		}
	}
}

func printLine(line string) {
	fmt.Printf("%s\n", line)
}

func main() {
	read(filename)
	os.Exit(found)
}
