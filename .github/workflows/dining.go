package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

const numPhilosophers = 6

type Chopstick struct{ sync.Mutex }

type Philosopher struct {
	id                  int
	leftCS, rightCS     *Chopstick
}

func (p Philosopher) dine(wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 0; i < 3; i++ { // each philosopher tries 3 meals
		fmt.Printf("[%d] Philosopher %d THINK\n", time.Now().UnixNano(), p.id)
		time.Sleep(time.Millisecond * time.Duration(rand.Intn(200)))

		// Odd-even strategy to avoid circular wait
		if p.id%2 == 0 {
			// Even philosopher:left then right
			fmt.Printf("[%d] Philosopher %d WAIT left\n", time.Now().UnixNano(), p.id)
			p.leftCS.Lock()
			fmt.Printf("[%d] Philosopher %d WAIT right\n", time.Now().UnixNano(), p.id)
			p.rightCS.Lock()
		} else {
			// Odd philosopher:right then left
			fmt.Printf("[%d] Philosopher %d WAIT right\n", time.Now().UnixNano(), p.id)
			p.rightCS.Lock()
			fmt.Printf("[%d] Philosopher %d WAIT left\n", time.Now().UnixNano(), p.id)
			p.leftCS.Lock()
		}

		// Eating
		fmt.Printf("[%d] Philosopher %d EAT\n", time.Now().UnixNano(), p.id)
		time.Sleep(time.Millisecond * time.Duration(rand.Intn(200)))

		// Release chopsticks
		p.leftCS.Unlock()
		p.rightCS.Unlock()
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())

	// Initialize chopsticks
	chopsticks := make([]*Chopstick, numPhilosophers)
	for i := 0; i < numPhilosophers; i++ {
		chopsticks[i] = new(Chopstick)
	}

	// Initialize philosophers
	philosophers := make([]*Philosopher, numPhilosophers)
	for i := 0; i < numPhilosophers; i++ {
		philosophers[i] = &Philosopher{
			id: i,
			leftCS: chopsticks[i],
			rightCS: chopsticks[(i+1)%numPhilosophers],
		}
	}

	// Start dining
	var wg sync.WaitGroup
	for i := 0; i < numPhilosophers; i++ {
		wg.Add(1)
		go philosophers[i].dine(&wg)
	}

	wg.Wait()
	fmt.Println("Simulation finished.")
}
