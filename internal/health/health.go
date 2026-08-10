// Package health は gRPC 標準の health checking protocol
// (grpc.health.v1.Health) を提供する。
//
// 目的は **プロセスの生存ではなく実際に応答できるかを外から判定できるようにする**こと。
// プロセスが生きたままハングしたり DB を見失ったりした場合、TCP の生存確認だけでは
// 検知できず、k8s の Pod は Ready のままトラフィックを受け続けてしまう。
package health

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

const (
	// DB への ping 間隔。k8s の probe 間隔より短くしておく。
	checkInterval = 5 * time.Second
	// ping 自体のタイムアウト。これを超えたら NOT_SERVING と判定する。
	checkTimeout = 2 * time.Second
)

// Checker は DB の疎通結果を gRPC health のステータスに反映する。
type Checker struct {
	server *health.Server
	pool   *pgxpool.Pool
}

func NewChecker(pool *pgxpool.Pool) *Checker {
	s := health.NewServer()
	// 起動直後はまだ疎通を確認していないので NOT_SERVING から始める。
	// Watch() している側に「準備中」を正しく伝えるため。
	s.SetServingStatus("", healthpb.HealthCheckResponse_NOT_SERVING)
	return &Checker{server: s, pool: pool}
}

// Register は health service を gRPC サーバに登録する。
func (c *Checker) Register(s *grpc.Server) {
	healthpb.RegisterHealthServer(s, c.server)
}

// Start は DB 疎通の監視を開始する。ctx が終わると NOT_SERVING にして戻る。
// 呼び出し側は goroutine で回すこと。
func (c *Checker) Start(ctx context.Context) {
	c.check(ctx)

	ticker := time.NewTicker(checkInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			// shutdown 時は NOT_SERVING を publish してから抜ける。
			// LB / k8s に「もう振らないで」を先に伝えるため。
			c.server.SetServingStatus("", healthpb.HealthCheckResponse_NOT_SERVING)
			return
		case <-ticker.C:
			c.check(ctx)
		}
	}
}

// check は DB に ping して結果をステータスへ反映する。
func (c *Checker) check(ctx context.Context) {
	pingCtx, cancel := context.WithTimeout(ctx, checkTimeout)
	defer cancel()

	if err := c.pool.Ping(pingCtx); err != nil {
		// ctx 自体が終了しているときは shutdown なのでエラー扱いしない。
		if ctx.Err() != nil {
			return
		}
		log.Ctx(ctx).Warn().Err(err).Msg("health check failed: database unreachable")
		c.server.SetServingStatus("", healthpb.HealthCheckResponse_NOT_SERVING)
		return
	}

	c.server.SetServingStatus("", healthpb.HealthCheckResponse_SERVING)
}
