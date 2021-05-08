/*
 Feedback Systems
 Chapter.1 - Buffer controller

*/

package main

import (
	"fmt"
	"math"
	"math/rand"

	v1 "github.com/timgluz/pidip/pkg/pidip"
)

// TODO: read configs from CLI
// TODO: add plotter with gnuplot
func main() {
	var steps uint = 1000
	var use_open = true

	controller := v1.NewPIController(1.25, 0.01)
	buffer := NewBuffer(50, 10)

	if use_open {
		run_open_loop(&buffer, steps)
	} else {
		run_closed_loop(&controller, &buffer, steps)
	}
}

func run_open_loop(buffer *Buffer, max_steps uint) {
	var targets []uint
	var processed_units []uint
	var e uint = 0

	for t := uint(0); t < max_steps; t++ {
		var u uint = setpoint(t)
		y := buffer.Work(u)
		targets = append(targets, u)
		processed_units = append(processed_units, y)

		fmt.Printf("step: %d, r: %d, e: %d, u: %d, y: %d\n", t, u, e, u, y)
	}
}

func run_closed_loop(controller *v1.PIController, buffer *Buffer, max_steps uint) {
	fmt.Println("TBD")
}

// returns simulation reference point
func setpoint(step uint) uint {
	switch {
	case step < 100:
		return 0
	case step < 300:
		return 50
	default:
		return 10
	}
}

/*
Example 1: Buffer


      ┌───────────┐                              ┌────────────────┐
      │           │     ┌─────────────────┐      │                    ┌────────────────┐
─────►│ Readypool ├────►│   Randomness    ├─────►│     Buffer     ├──►│   Randomness   │
      │           │     └─────────────────┘      │                │   └────────────────┘
      └───────────┘                              └────────────────┘
*/

type Buffer struct {
	max_wip  uint
	max_flow uint
	queued   uint
	wip      uint // work in progress
}

func NewBuffer(max_wip, max_flow uint) Buffer {
	return Buffer{max_wip, max_flow, 0, 0}
}

func (b *Buffer) Work(u uint) uint {
	// add random amount of u to ready pool
	u_rand := rand.Float64() * float64(u)
	u_low := math.Min(0.0, u_rand)
	u_max := math.Max(u_low, float64(b.max_wip))

	b.wip = uint(math.Floor(u_max))

	// move from ready pool to queue
	r := uint(rand.Intn(int(b.max_flow)))
	if r > b.wip {
		r = b.wip
	}

	b.wip -= r
	b.queued += r

	// release from queue to downstream process
	done := uint(rand.Intn(int(b.max_flow)))
	if done > b.queued {
		done = b.queued
	}
	b.queued -= done

	return b.queued
}
