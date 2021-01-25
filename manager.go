package GoTrivia

import (
	"fmt"
	"log"
	"math/rand"
	"strings"
	"sync"
	"time"
)

type Manager struct {
	Players     map[string]*Player
	mutex       sync.Mutex
	currentGame *Game
	out         chan Message
}

type Message struct {
	From *Player
	Message string
}

func NewManager() *Manager {
	return &Manager{
		mutex:   sync.Mutex{},
		Players: make(map[string]*Player, 0),
		out:     make(chan Message),
	}
}

func (tg *Manager) NewGame() {
	tg.mutex.Lock()
	defer tg.mutex.Unlock()

	now := time.Now()
	s := rand.NewSource(time.Now().Unix())
	r := rand.New(s)

	tg.currentGame = &Game{
		Id:        now.UnixNano(),
		Question:  questions[r.Intn(len(questions))],
		StartTime: time.Now(),
	}
	log.Printf("Question: %s", tg.currentGame.Question.question)
	log.Printf("Answer: %s", tg.currentGame.Question.answer)

	tg.Send(fmt.Sprintf("Question: %s?", tg.currentGame.Question.question), nil)
}

func (tg *Manager) CurrentGame() *Game {
	tg.mutex.Lock()
	defer tg.mutex.Unlock()

	return tg.currentGame
}

func (tg *Manager) TryAnswer(p *Player, a string) {
	tg.mutex.Lock()

	if !strings.EqualFold(tg.currentGame.Question.answer, a) {
		tg.Send(fmt.Sprintf("..::Wrong Answer by: %s::..", p.Name), p)
		tg.mutex.Unlock()
		return
	}

	if player, found := tg.Players[p.Id]; found {
		player.Correct++
	} else {
		p.Correct = 1
		tg.Players[p.Id] = p
	}

	tg.Send(fmt.Sprintf("..::Correct Answer By: %s::..", p.Name), p)

	tg.mutex.Unlock()
	tg.NewGame()
}

func (tg *Manager) GetQuestion() {
	tg.mutex.Lock()
	defer tg.mutex.Unlock()
	tg.Send(fmt.Sprintf("..::Current Question: %s? ::..", tg.currentGame.Question.question), nil)
}

func (tg *Manager) GetScore(p *Player) {
	tg.mutex.Lock()
	defer tg.mutex.Unlock()
	if player, found := tg.Players[p.Id]; found {
		tg.Send(fmt.Sprintf("...::%s Score: %d", p.Name, player.Correct), nil)
	} else {
		tg.Send(fmt.Sprintf("...::%s Score: 0", p.Name), nil)
	}
}
// communication
func (tg *Manager) Send(m string, p *Player) {
	tg.Outgoing() <- Message{
		Message: m,
		From:    p,
	}
}

func (tg *Manager) Outgoing() chan Message {
	return tg.out
}
