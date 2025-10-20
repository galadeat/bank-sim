package repl

import (
	"bufio"
	"fmt"
	"strings"
)

// readInput prompts and returns a trimmed line.
func readInput(reader *bufio.Reader, prompt string) string {
	fmt.Print(prompt)
	line, _ := reader.ReadString('\n')
	return strings.TrimSpace(line)
}
