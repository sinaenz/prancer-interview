package service

import (
	"context"
	"fmt"
	"math"
	"time"
)

type Coordinator interface {
	SelectAgent(ctx context.Context, x, y int) (agentNum int, err error)
	MoveAgent(ctx context.Context, x, y int)
}

type target struct {
	X float64
	Y float64
}

type agent struct {
	ID       int
	X        float64
	Y        float64
	Moving   bool
	Distance float64
	Target   *target
	Command  chan interface{}
}

func (a *agent) Listen(t *target) {
	switch cmd := <-a.Command; cmd.(type) {
	case *target:
		a.Target = cmd.(*target)
		a.Distance = math.Sqrt(math.Pow(cmd.(*target).X-a.X, 2) + math.Pow(cmd.(*target).Y-a.Y, 2))
		a.Command <- a.Distance

	case bool:
		if a.Moving = cmd.(bool); a.Moving {
			// Lock command channel
			a.Command <- struct{}{}
			for i := 0; (a.Distance) >= 1; i++ {
				a.Distance -= 1
				fmt.Println(fmt.Sprintf("Agent %d has %f unit to destination", a.ID, a.Distance))
				time.Sleep(time.Second)
			}
			fmt.Println(fmt.Sprintf("Agent %d has %f unit to destination", a.ID, a.Distance))
			time.Sleep(time.Duration(int(a.Distance*1000)) * time.Millisecond)
		}
		a.Moving = false
		a.Target = nil
		// Release Command Channel
		<-a.Command
	}
}

type center struct {
	Agents []*agent
}

func NewCenter() *center {
	var agents []*agent
	for i := 0; i < 8; i++ {
		agents = append(agents, &agent{
			ID:     i,
			X:      .0,
			Y:      .0,
			Moving: false,
		})
	}
	return &center{
		Agents: agents,
	}
}

func (c *center) SelectAgent(ctx context.Context, x, y float64) (agentNum int, err error) {
	panic("Imp")
}

func (c *center) MoveAgent(ctx context.Context, x, y float64) {
	//TODO implement me
	panic("implement me")
}

func (c *center) findNearestAgent(ctx context.Context, x, y float64) (agentID int, agentDistance float64) {
	distance := .0
	for i, a := range c.Agents {
		if a.Moving {
			continue
		}
		d := math.Sqrt(math.Pow(x-a.X, 2) + math.Pow(y-a.Y, 2))
		if d < distance {
			agentID = i
			agentDistance = d
		}
	}
	return agentID, agentDistance
}
