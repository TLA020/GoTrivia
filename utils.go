package GoTrivia

import (
	"bufio"
	"log"
	"os"
	"runtime"
	"strings"
)

func scanLinesToArray(p string) (r []string, b bool) {
	_, f, _, _ := runtime.Caller(0)
	lastOccurrence := strings.LastIndex(f, "/")
	dir := f[0:lastOccurrence]

	file, err := os.Open(dir + p)
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


