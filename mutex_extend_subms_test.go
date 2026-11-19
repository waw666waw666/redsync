package redsync

import (
	"testing"
	"time"
)

// Extend used to pass int(expiry/time.Millisecond) to PEXPIRE. A positive TTL
// below 1ms became 0, so Redis deleted the key. Lock with a normal expiry
// first, then shrink m.expiry before Extend so every driver hits that path.
func TestExtendSubMillisecondKeepsLock(t *testing.T) {
	for k, v := range makeCases(1) {
		t.Run(k, func(t *testing.T) {
			rs := New(v.pools...)
			mutex := rs.NewMutex("test-subms-extend", WithExpiry(time.Second))
			if err := mutex.Lock(); err != nil {
				t.Fatalf("Lock: %v", err)
			}
			mutex.expiry = 500 * time.Microsecond
			// Extend local until is ~1ms after rounding; slow drivers may still
			// return ErrExtendFailed. The bug was deleting the Redis key.
			_, _ = mutex.Extend()
			valid, _ := mutex.Valid()
			if !valid {
				t.Fatal("lock key missing after Extend with sub-millisecond expiry")
			}
			if _, err := mutex.Unlock(); err != nil {
				t.Fatalf("Unlock: %v", err)
			}
		})
	}
}
