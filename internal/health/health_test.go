package health

import (
	"context"
	"errors"
	"testing"
	"time"

	"google.golang.org/grpc/health/grpc_health_v1"
)

type fakePinger struct{ err error }

func (f *fakePinger) Ping(context.Context) error { return f.err }

func currentStatus(t *testing.T, c *Checker) grpc_health_v1.HealthCheckResponse_ServingStatus {
	t.Helper()
	resp, err := c.Check(context.Background(), &grpc_health_v1.HealthCheckRequest{})
	if err != nil {
		t.Fatalf("Check: %v", err)
	}
	return resp.GetStatus()
}

// Check が直近の ping 結果を反映することを DB 無しで検証する。
func TestCheckReflectsPing(t *testing.T) {
	fp := &fakePinger{}
	c := &Checker{db: fp}
	c.set(grpc_health_v1.HealthCheckResponse_NOT_SERVING)

	if got := currentStatus(t, c); got != grpc_health_v1.HealthCheckResponse_NOT_SERVING {
		t.Fatalf("initial: want NOT_SERVING, got %v", got)
	}

	fp.err = nil
	c.check(context.Background())
	if got := currentStatus(t, c); got != grpc_health_v1.HealthCheckResponse_SERVING {
		t.Fatalf("after ok ping: want SERVING, got %v", got)
	}

	fp.err = errors.New("db down")
	c.check(context.Background())
	if got := currentStatus(t, c); got != grpc_health_v1.HealthCheckResponse_NOT_SERVING {
		t.Fatalf("after failed ping: want NOT_SERVING, got %v", got)
	}
}

// set は遷移したときだけ changed=true を返す(遷移時のみログするため)。
func TestSetReportsTransition(t *testing.T) {
	c := &Checker{db: &fakePinger{}}
	c.set(grpc_health_v1.HealthCheckResponse_NOT_SERVING)

	if !c.set(grpc_health_v1.HealthCheckResponse_SERVING) {
		t.Fatal("NOT_SERVING -> SERVING should report changed")
	}
	if c.set(grpc_health_v1.HealthCheckResponse_SERVING) {
		t.Fatal("SERVING -> SERVING should report no change")
	}
}

// Start は ctx cancel で NOT_SERVING にして戻る。
func TestStartStopsOnCancel(t *testing.T) {
	c := &Checker{db: &fakePinger{}}
	ctx, cancel := context.WithCancel(context.Background())

	done := make(chan struct{})
	go func() {
		c.Start(ctx)
		close(done)
	}()

	cancel()

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("Start did not return after ctx cancel")
	}

	if got := currentStatus(t, c); got != grpc_health_v1.HealthCheckResponse_NOT_SERVING {
		t.Fatalf("after cancel: want NOT_SERVING, got %v", got)
	}
}
