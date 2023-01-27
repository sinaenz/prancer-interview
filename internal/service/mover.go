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
	Command  chan interface{}
	Response chan float64
}

func (a *agent) Listen() {
	for {
		select {
		case t := <-a.Command:
			switch t.(type) {
			case *Target: // set Target
				a.Target = t.(*Target)
				a.Distance = math.Sqrt(math.Pow(t.(*Target).X-a.X, 2) + math.Pow(t.(*Target).Y-a.Y, 2))
				a.Response <- a.Distance

			case bool: // start moving
				if t.(bool) {
					for i := 0; (a.Distance) > 1; i++ {
						fmt.Println(fmt.Sprintf("Agent %d has %f unit to destination", a.ID, a.Distance))
						time.Sleep(time.Second)
						a.Distance -= 1
					}
					fmt.Println(fmt.Sprintf("Agent %d has %f unit to destination", a.ID, a.Distance))
					time.Sleep(time.Duration(int(a.Distance*1000)) * time.Millisecond)
					fmt.Println(fmt.Sprintf("Agent %d has received!", a.ID))

				}
				a.X = a.Target.X
				a.Y = a.Target.Y
				a.Target = nil
			}
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
		a := agent{
			ID:       i,
			X:        .0,
			Y:        .0,
			Distance: 0,
			Target:   nil,
			Command:  make(chan interface{}),
			Response: make(chan float64),
		}
		agents = append(agents, &a)
		go a.Listen()
	}
	return &center{
		Agents: agents,
	}
}

func (c *center) SelectAgent(ctx context.Context, t *Target) (agent *agent, err error) {
	distance := math.MaxFloat64
	for i, a := range c.Agents {
		select {
		case a.Command <- t:
			d := <-a.Response
			if d < distance {
				agent = a
				distance = d
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
	if err != nil {
		return err
	}
	a.Command <- true
	return err
}
