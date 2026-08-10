// Package gateway は API ゲートウェイ(ESP/Envoy 等)が付与する検証済みユーザー情報ヘッダを
// 解釈し subject(識別子)を取り出す。gRPC には依存しない（ヘッダ値の解釈ロジックのみ）。
//
// 前提: app はゲートウェイ経由でのみ到達可能（loopback / 内部ネットワーク限定で公開）。
// 直接到達できるとヘッダを偽造されて認証を素通りできるため、必ず前段を挟むこと。
package gateway

import (
	"encoding/base64"
	"encoding/json"
	"errors"
)

// HeaderUserInfo は Cloud Endpoints/ESP(および同等の Envoy jwt_authn 設定)が検証済み
// JWT payload を base64url(JSON)で載せて転送するヘッダ名。
// https://cloud.google.com/endpoints/docs/grpc/authenticating-users
const HeaderUserInfo = "x-endpoint-api-userinfo"

// ErrNoSubject はヘッダに sub が無い場合。
var ErrNoSubject = errors.New("user info has no subject")

// userInfo は HeaderUserInfo の中身のうち必要な部分。
// sub は JWT の標準 subject クレーム（IdP 非依存の不透明な識別子）。
type userInfo struct {
	Subject string `json:"sub"`
}

// Parse は HeaderUserInfo の値(base64url(JSON))から subject(sub)を取り出す。
// sub は IdP が発行する不透明な文字列で、app 内ではそのまま user id として扱う
// （ULID とは限らない。例: Google OAuth は 21 桁の数値文字列）。
func Parse(encoded string) (string, error) {
	raw, err := base64.RawURLEncoding.DecodeString(encoded)
	if err != nil {
		// padding 付きのケースにもフォールバック
		raw, err = base64.StdEncoding.DecodeString(encoded)
		if err != nil {
			return "", err
		}
	}

	var info userInfo
	if err := json.Unmarshal(raw, &info); err != nil {
		return "", err
	}

	if info.Subject == "" {
		return "", ErrNoSubject
	}

	return info.Subject, nil
}
