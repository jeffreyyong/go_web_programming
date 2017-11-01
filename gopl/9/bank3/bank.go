// Package bank provides a concurrency-safe single-account bank.

package bank

import "sync"

var (
	mu      sync.Mutex // guards balance
	balance int
)

func Deposit(amount int) {
	mu.Lock()
	defer mu.Unlock()
	deposit(amount)
}

func Balance() int {
	mu.Lock()
	mu.Unlock()
	return balance
}

// Deposit tries to acquire the mutex lock a secondtime by calling mu.Lock(), but because mutex locks
// are not re-entrant - it's not possible to lock a mutex that's already locked - this leads to a deadlock
// where nothing can proceed, and Withdraw blocks forever

func Withdraw(amount int) bool {
	mu.Lock()
	defer mu.Unlock()
	deposit(-amount)
	if balance < 0 {
		deposit(amount)
		return false // insufficient funds
	}
	return true
}

// This function requires that the lock be held.
func deposit(amount int) { balance += amount }
