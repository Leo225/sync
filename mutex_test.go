package sync

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
)

func TestRedisMutex(t *testing.T) {
	var ctx = context.Background()
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	pong, err := rdb.Ping(ctx).Result()
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(pong)

	rs := NewRedisSync(rdb)
	mux := rs.NewMutex("test_key")
	err = mux.Lock()
	if err != nil {
		t.Errorf("Lock: %s", err)
	}
	for i := 0; i < 10; i++ {
		err = mux.Lock()
		if err != nil {
			t.Errorf("Lock: %s", err)
		}
		time.Sleep(2 * time.Second)
		if i%2 == 0 {
			fmt.Println("double multiple unlock")
			err = mux.Unlock()
			if err != nil {
				t.Errorf("Unlock: %s", err)
			}
		}
	}
	err = mux.Unlock()
	if err != nil {
		t.Errorf("Unlock: %s", err)
	}
}

func BenchmarkRedisMutex(b *testing.B) {
	var ctx = context.Background()
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	pong, err := rdb.Ping(ctx).Result()
	if err != nil {
		b.Error(err)
		return
	}
	b.Log(pong)
	rs := NewRedisSync(rdb)
	b.ResetTimer()
	mux := rs.NewMutex("test_benchmark")
	for i := 0; i < b.N; i++ {
		err = mux.Lock()
		if err != nil {
			b.Errorf("Lock: %s", err)
		}
		err = mux.Unlock()
		if err != nil {
			b.Errorf("Unlock: %s", err)
		}
	}
}
