package GoTrivia

import (
	"bufio"
	"log"
	"os"
)

func scanLinesToArray(p string) (r []string, b bool) {
	file, err := os.Open(p)
	if err != nil {
		log.Print(err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		r = append(r, line)
	}
	return r, true
}


