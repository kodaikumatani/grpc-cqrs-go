package gateway

import (
	"encoding/base64"
	"testing"
)

func encode(t *testing.T, jsonBody string) string {
	t.Helper()
	return base64.RawURLEncoding.EncodeToString([]byte(jsonBody))
}

func TestParseUserID(t *testing.T) {
	tests := []struct {
		name    string
		encoded string
		want    string
		wantErr bool
	}{
		{
			name:    "sub から取得(Google の数値 sub)",
			encoded: encode(t, `{"sub":"117012345678901234567","email":"a@example.com"}`),
			want:    "117012345678901234567",
		},
		{
			name:    "sub は不透明文字列でよい(ULID でなくても可)",
			encoded: encode(t, `{"sub":"auth0|abc123"}`),
			want:    "auth0|abc123",
		},
		{
			name:    "sub 空は失敗",
			encoded: encode(t, `{"email":"a@example.com"}`),
			wantErr: true,
		},
		{
			name:    "base64 でない",
			encoded: "%%%not-base64%%%",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse(tt.encoded)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("want error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Fatalf("want %s, got %s", tt.want, got)
			}
		})
	}
}
