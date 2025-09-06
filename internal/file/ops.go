package file

import (
	"bufio"
	"os"
	"strings"
)

// Load reads a file and returns its content as a slice of rune slices.
func Load(fileName string) ([][]rune, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var content [][]rune
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		content = append(content, []rune(line))
	}

	return content, scanner.Err()
}

// Save writes the editor's content back to a file.
func Save(fileName string, content [][]rune) error {
	var lines []string
	for _, line := range content {
		lines = append(lines, string(line))
	}
	return os.WriteFile(fileName, []byte(strings.Join(lines, "\n")), 0644)
}
