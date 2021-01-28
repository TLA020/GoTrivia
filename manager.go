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

func (m *Manager) NewGame() {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	now := time.Now()
	s := rand.NewSource(time.Now().Unix())
	r := rand.New(s)

	m.currentGame = &Game{
		Id:        now.UnixNano(),
		Question:  questions[r.Intn(len(questions))],
		StartTime: time.Now(),
	}
	log.Printf("[Question]: %s", m.currentGame.Question.question)
	log.Printf("[Answer]: %s", m.currentGame.Question.answer)

	m.Send(fmt.Sprintf("[New Question]: %s?::..", m.currentGame.Question.question), nil)
}

func (m *Manager) CurrentGame() *Game{
	m.mutex.Lock()
	defer m.mutex.Unlock()

	return m.currentGame
}

func (m *Manager) TryAnswer(p *Player, a string) {
	if m.currentGame == nil {
		return
	}

	m.mutex.Lock()

	if !strings.EqualFold(m.currentGame.Question.answer, a) {
		m.Send(fmt.Sprintf("[Wrong Answer]: %s", p.Name), p)
		m.Send(fmt.Sprintf("[Question]: %s?", m.currentGame.Question.question), p)
		m.mutex.Unlock()
		return
	}

	if player, found := m.Players[p.Id]; found {
		player.Correct++
	} else {
		p.Correct = 1
		m.Players[p.Id] = p
	}

	m.Send(fmt.Sprintf("[Correct Answer]: %s", p.Name), p)

	m.mutex.Unlock()
	m.NewGame()
}

func (m *Manager) GetQuestion() {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if m.currentGame == nil {
		m.Send("[NO GAME]", nil)
	} else {
		m.Send(fmt.Sprintf("[Current Question]: %s? ", m.currentGame.Question.question), nil)
	}
}

func (m *Manager) Skip() {
	m.mutex.Lock()
	if m.currentGame != nil {
		m.Send(fmt.Sprintf("[Skipped] Answer: %s ", m.currentGame.Question.answer), nil)
	}
	m.mutex.Unlock()
	m.NewGame()
}

func (m *Manager) GetScore(p *Player) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	if player, found := m.Players[p.Id]; found {
		m.Send(fmt.Sprintf("[%s Score]: %d", p.Name, player.Correct), nil)
	} else {
		m.Send(fmt.Sprintf("[%s Score]: 0", p.Name), nil)
	}
}

// communication
func (m *Manager) Send(msg string, p *Player) {
	m.Outgoing() <- Message{
		Message: msg,
		From:    p,
	}
}

func (m *Manager) Outgoing() chan Message {
	return m.out
}
