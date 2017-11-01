// Package bank provides a concurrency-safe bank with one account
package bank1

var deposits = make(chan int) // send amount to deposit
var balances = make(chan int) // receive balance

func Deposit(amount int) { deposits <- amount }
func Balance() int       { return <-balances }

func teller() {
	var balance int // balance is confined to teller goroutine

	for {
		select {
		case amount := <-deposits:
			balance += amount
		case balances <- balance:
		}
	}
}

func init() {
	go teller() // Start the monitor goroutine
}

// Even when a variable cannot be confined to a single goroutine for its entire lifetime confiemnet
// may still be a solution to the problem of concurrent access. For example, it's common to share a variable between
// goroutines in a pipeline by passing its address from one stage to next over a channel
