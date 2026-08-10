// Package gateway は API ゲートウェイ(ESP 等)が付与する検証済みユーザー情報ヘッダを
// 解釈し UID を取り出す。gRPC には依存しない（ヘッダ値の解釈ロジックのみ）。
//
// 前提: app はゲートウェイ経由でのみ到達可能（loopback / 内部ネットワーク限定で公開）。
// 直接到達できるとヘッダを偽造されて認証を素通りできるため、必ず前段を挟むこと。
package gateway

import (
	"encoding/base64"
	"encoding/json"

	"github.com/oklog/ulid/v2"
)

// HeaderUserInfo は Cloud Endpoints/ESP が検証済み JWT payload を base64url(JSON)
// で載せて転送するヘッダ名。ESP の固定仕様。
// https://cloud.google.com/endpoints/docs/grpc/authenticating-users
const HeaderUserInfo = "x-endpoint-api-userinfo"

// userInfo は HeaderUserInfo の中身のうち必要な部分。
// sub は JWT の標準 subject クレーム（IdP 非依存）。Firebase の uid もここに入る。
type userInfo struct {
	Subject string `json:"sub"`
}

// Parse は HeaderUserInfo の値(base64url(JSON))から UID を取り出す。
func Parse(encoded string) (ulid.ULID, error) {
	raw, err := base64.RawURLEncoding.DecodeString(encoded)
	if err != nil {
		// padding 付きのケースにもフォールバック
		raw, err = base64.StdEncoding.DecodeString(encoded)
		if err != nil {
			return ulid.ULID{}, err
		}
	}

	var info userInfo
	if err := json.Unmarshal(raw, &info); err != nil {
		return ulid.ULID{}, err
	}

	return ulid.Parse(info.Subject)
}
