package main

import (
	"flag"
	"fmt"
	"strings"
)

var (
	from, to      string
	limit, offset int64
)

func init() {
	flag.StringVar(&from, "from", "", "file to read from")
	flag.StringVar(&to, "to", "", "file to write to")
	flag.Int64Var(&limit, "limit", 0, "limit of bytes to copy")
	flag.Int64Var(&offset, "offset", 0, "offset in input file")
}

func main() {
	flag.Parse()

	loadingCh := make(chan int64)

	go func() {
		if err := Copy(from, to, offset, limit, loadingCh); err != nil {
			fmt.Println("Error: ", err)
		}
	}()

	for c := range loadingCh {
		fmt.Print(getProgressbar(c))
	}

	fmt.Print("\n")
}

func getProgressbar(percent int64) string {
	var result strings.Builder
	result.WriteString("\r\033[0K[")

	maxBarLength := int64(20)
	filledBarLength := percent / (100 / maxBarLength)

	for range filledBarLength {
		result.WriteString("#")
	}
	for range maxBarLength - filledBarLength {
		result.WriteString(" ")
	}

	result.WriteString(fmt.Sprintf("] %d%%", percent))
	return result.String()
}
