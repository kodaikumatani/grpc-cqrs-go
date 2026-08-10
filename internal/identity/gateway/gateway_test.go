package gateway

import (
	"encoding/base64"
	"testing"

	"github.com/oklog/ulid/v2"
)

func encode(t *testing.T, jsonBody string) string {
	t.Helper()
	return base64.RawURLEncoding.EncodeToString([]byte(jsonBody))
}

func TestParseUserID(t *testing.T) {
	uid := ulid.Make().String()

	tests := []struct {
		name    string
		encoded string
		wantID  string
		wantErr bool
	}{
		{
			name:    "sub から取得",
			encoded: encode(t, `{"sub":"`+uid+`","email":"a@example.com"}`),
			wantID:  uid,
		},
		{
			name:    "非 ULID は失敗",
			encoded: encode(t, `{"sub":"not-a-ulid"}`),
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
			if got.String() != tt.wantID {
				t.Fatalf("want %s, got %s", tt.wantID, got.String())
			}
		})
	}
}
