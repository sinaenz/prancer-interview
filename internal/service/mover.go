package service

import (
	"context"
	"errors"
	"fmt"
	"math"
	"time"
)

type Target struct {
	X float64
	Y float64
}

type agent struct {
	ID       int
	X        float64
	Y        float64
	Distance float64
	Target   *Target
	Lock     chan bool
	Command  chan interface{}
	Response chan float64
	Start    chan bool
}

func (a *agent) Listen() {
	select {
	case t := <-a.Command:
		switch t.(type) {
		case *Target: // set Target
			a.Target = t.(*Target)
			a.Distance = math.Sqrt(math.Pow(t.(*Target).X-a.X, 2) + math.Pow(t.(*Target).Y-a.Y, 2))
			a.Response <- a.Distance

		case bool: // start moving
			if t.(bool) {
				for i := 0; (a.Distance) >= 1; i++ {
					a.Distance -= 1
					fmt.Println(fmt.Sprintf("Agent %d has %f unit to destination", a.ID, a.Distance))
					time.Sleep(time.Second)
				}
				fmt.Println(fmt.Sprintf("Agent %d has %f unit to destination", a.ID, a.Distance))
				time.Sleep(time.Duration(int(a.Distance*1000)) * time.Millisecond)
			}
			a.Lock <- false
			a.Target = nil
		}
	}
}

type Center interface {
	SelectAgent(ctx context.Context, target *Target) (agent *agent, err error)
	MoveAgent(ctx context.Context, target *Target) (err error)
}

type center struct {
	Agents []*agent
}

func NewCenter() Center {
	var agents []*agent
	for i := 0; i < 8; i++ {
		agents = append(agents, &agent{
			ID:       i,
			X:        .0,
			Y:        .0,
			Distance: 0,
			Target:   nil,
			Lock:     make(chan bool),
			Command:  make(chan interface{}),
			Response: make(chan float64),
			Start:    nil,
		})
	}
	return &center{
		Agents: agents,
	}
}

func (c *center) SelectAgent(ctx context.Context, t *Target) (agent *agent, err error) {
	distance := .0
	for i, a := range c.Agents {
		select {
		case a.Lock <- true:
			a.Command <- t
			d := <-a.Response
			if d < distance {
				agent = a
			}
		default:
			fmt.Println(fmt.Sprintf("Agent %d is busy!", i))
		}
	}
	if agent == nil {
		err = errors.New("no agent is available")
	}
	return agent, err
}

func (c *center) MoveAgent(ctx context.Context, t *Target) (err error) {
	a, err := c.SelectAgent(ctx, t)
	a.Start <- true
	return err
}
