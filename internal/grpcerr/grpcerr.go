// Package grpcerr は client 向けの安全な gRPC status と、ログ用の元 cause を
// 同時に運ぶためのエラー型を提供する。特定の app/domain には依存しない。
package grpcerr

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// statusError は GRPCStatus() を直接実装することで client には安全な status を
// 返しつつ、Error()/Unwrap() で元の cause を保持しログに詳細を残せるようにする。
//
// 注意: status error を fmt.Errorf("%w") で包むと status.FromError が message を
// err.Error() 全文に差し替えて client に漏れる。必ずこの型（GRPCStatus 直接実装）
// を使い、%w ラップはしないこと。
type statusError struct {
	cause error
	st    *status.Status
}

func (e *statusError) Error() string              { return e.cause.Error() }
func (e *statusError) Unwrap() error              { return e.cause }
func (e *statusError) GRPCStatus() *status.Status { return e.st }

// WithStatus は cause を保持したまま client 向けの安全な status(code, msg) を付与する。
func WithStatus(cause error, code codes.Code, msg string) error {
	return &statusError{
		cause: cause,
		st:    status.New(code, msg),
	}
}
