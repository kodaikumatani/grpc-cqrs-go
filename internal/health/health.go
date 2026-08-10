// Package health は gRPC 標準の health checking protocol
// (grpc.health.v1.Health) を提供する。
//
// 目的は **プロセスの生存ではなく実際に応答できるかを外から判定できるようにする**こと。
// プロセスが生きたままハングしたり DB を見失ったりした場合、TCP の生存確認だけでは
// 検知できず、k8s の Pod は Ready のままトラフィックを受け続けてしまう。
//
// ⚠️ この service を実際に参照するのは前段の Envoy(app cluster の grpc_health_check)。
// Envoy 側が tcp から grpc に切り替わって初めて readiness に反映される。
//
// google.golang.org/grpc/health の helper サーバは使わず、grpc_health_v1 の
// 生成インターフェースを直接実装する(必要なのは Check の unary だけ)。
package health

import (
	"context"
	"sync/atomic"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
)

const (
	// DB への ping 間隔。k8s の probe 間隔より短くしておく。
	checkInterval = 5 * time.Second
	// ping 自体のタイムアウト。これを超えたら NOT_SERVING と判定する。
	checkTimeout = 2 * time.Second
)

// pinger は疎通確認に必要な最小インターフェース(*pgxpool.Pool が満たす)。テスト差し替え用。
type pinger interface {
	Ping(ctx context.Context) error
}

// Checker は grpc.health.v1.Health を実装し、DB の疎通結果を Check のステータスに反映する。
type Checker struct {
	grpc_health_v1.UnimplementedHealthServer
	db     pinger
	status atomic.Int32 // grpc_health_v1.HealthCheckResponse_ServingStatus
}

func NewChecker(pool *pgxpool.Pool) *Checker {
	c := &Checker{db: pool}
	// 起動直後はまだ疎通を確認していないので NOT_SERVING から始める。
	c.set(grpc_health_v1.HealthCheckResponse_NOT_SERVING)
	return c
}

// Register は health service を gRPC サーバに登録する。
func (c *Checker) Register(s *grpc.Server) {
	grpc_health_v1.RegisterHealthServer(s, c)
}

// Check は現在の serving status を返す(unary)。Envoy / kubelet の gRPC ヘルスチェックが使う。
// service 名は問わず全体のステータスを返す(単一サービス構成のため)。
// Watch(stream) は未対応 = UnimplementedHealthServer に委ねる(Unimplemented)。
func (c *Checker) Check(
	_ context.Context,
	_ *grpc_health_v1.HealthCheckRequest,
) (*grpc_health_v1.HealthCheckResponse, error) {
	return &grpc_health_v1.HealthCheckResponse{
		Status: grpc_health_v1.HealthCheckResponse_ServingStatus(c.status.Load()),
	}, nil
}

// Start は DB 疎通の監視を開始する。ctx が終わると NOT_SERVING にして戻る。
// 呼び出し側は goroutine で回すこと。
// なお ctx cancel 時の NOT_SERVING が LB/k8s のドレインとして効くのは、
// SIGTERM を捕まえて GracefulStop する場合のみ(現状はシグナル未処理なので効果は限定的)。
func (c *Checker) Start(ctx context.Context) {
	c.check(ctx)

	ticker := time.NewTicker(checkInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			c.set(grpc_health_v1.HealthCheckResponse_NOT_SERVING)
			return
		case <-ticker.C:
			c.check(ctx)
		}
	}
}

// check は DB に ping して結果をステータスへ反映する。ログは状態が遷移したときだけ出す
// (障害が続く間 5s ごとに出し続けるとログが溢れるため)。
func (c *Checker) check(ctx context.Context) {
	pingCtx, cancel := context.WithTimeout(ctx, checkTimeout)
	defer cancel()

	if err := c.db.Ping(pingCtx); err != nil {
		// ctx 自体が終了しているときは shutdown なのでエラー扱いしない。
		if ctx.Err() != nil {
			return
		}
		if c.set(grpc_health_v1.HealthCheckResponse_NOT_SERVING) {
			log.Ctx(ctx).Warn().Err(err).Msg("health: not serving (database unreachable)")
		}
		return
	}

	if c.set(grpc_health_v1.HealthCheckResponse_SERVING) {
		log.Ctx(ctx).Info().Msg("health: serving (database reachable)")
	}
}

// set はステータスを更新し、実際に変化したか(遷移か)を返す。
func (c *Checker) set(s grpc_health_v1.HealthCheckResponse_ServingStatus) (changed bool) {
	return c.status.Swap(int32(s)) != int32(s)
}
