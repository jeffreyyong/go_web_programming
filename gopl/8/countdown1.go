// Countdown implements the countdown for a rocket launch

package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println("Commencing countdown.")
	// time.Tick function returns a channel on which it sends events periodically, acting like a metronome
	tick := time.Tick(1 * time.Second)
	for countdown := 10; countdown > 0; countdown-- {
		fmt.Println(countdown)
		<-tick
	}
	launch()
}

func launch() {
	fmt.Println("Lift off!")
}
