package GoTrivia

import "log"

var questions []*Question

func init() {
	log.Print("init")
	q, ok := scanLinesToArray("./questions.txt")
	if !ok {
		return
	}

	a, ok := scanLinesToArray("./answers.txt")
	if !ok {
		return
	}

	if len(q) != len(a) {
		log.Print("Trivia: Question/Answer len not equal")
		return
	}

	for i, question := range q {
		questions = append(questions, &Question{
			question: question,
			answer: a[i],
		})
	}

	log.Printf("..::Trivia: loaded %d questions::..", len(q))
}
