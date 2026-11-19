package redsync

import (
	"testing"
	"time"
)

// WithExpiry below 1ms used to make Extend send PEXPIRE 0, which deletes the
// lock key while still returning success. go-redis truncates the initial
// SetNX up to 1ms, so Lock can succeed and then Extend would drop the key.
func TestExtendSubMillisecondKeepsLock(t *testing.T) {
	for k, v := range makeCases(1) {
		t.Run(k, func(t *testing.T) {
			rs := New(v.pools...)
			mutex := rs.NewMutex("test-subms-extend", WithExpiry(500*time.Microsecond))
			if err := mutex.Lock(); err != nil {
				t.Fatalf("Lock: %v", err)
			}
			ok, err := mutex.Extend()
			if err != nil {
				t.Fatalf("Extend: %v", err)
			}
			if !ok {
				t.Fatal("Extend returned false")
			}
			valid, err := mutex.Valid()
			if err != nil {
				t.Fatalf("Valid: %v", err)
			}
			if !valid {
				t.Fatal("lock key missing after Extend with sub-millisecond expiry")
			}
			if _, err := mutex.Unlock(); err != nil {
				t.Fatalf("Unlock: %v", err)
			}
		})
	}
}
