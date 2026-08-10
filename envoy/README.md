# Envoy ゲートウェイ (ESP の代替)

アプリの前段に Envoy を置き、JWT を検証して検証済みユーザー情報を
`X-Endpoint-API-UserInfo` ヘッダで backend に渡す構成。**GCP 非依存 / arm64 対応**なので
self-hosted k8s(Talos/RPi5 等)向け。ESP(Cloud Endpoints)を使う場合は `../espv2/` を参照。

## 仕組み

```
client ─(authorization: Bearer <JWT>)─▶ Envoy :8080
       ─(x-endpoint-api-userinfo: base64url(payload))─▶ app 127.0.0.1:50051
```

- **認証**: Envoy `jwt_authn` が JWKS で署名検証。`forward_payload_header` で payload を
  base64url(JSON) にしてヘッダへ → ESP の `X-Endpoint-API-UserInfo` と互換。
- **app は無改修**: `internal/identity/gateway.Parse` がこのヘッダから `sub` を読む。
- **認可(ReBAC)** は app 側(`internal/authz`)に残る。

## セキュリティ前提(重要)

app はヘッダを**信頼する**ので、**Envoy 経由でのみ到達可能**にすること。
本 scaffold は app を `-host 127.0.0.1`(loopback)で起動させ、Envoy を同一 Pod の
sidecar に置く前提。Service は Envoy の `:8080` だけを晒す。直接叩けるとヘッダ偽造で認証を破れる。

## 設定

`envoy.yaml` の以下を置換:

| プレースホルダ | 例(汎用 OIDC) |
|---|---|
| `<ISSUER>` | `https://your-idp.example.com/` |
| `<AUDIENCE>` | アプリの client_id / audience |
| `<JWKS_URI>` | `https://your-idp.example.com/.well-known/jwks.json` |
| `<JWKS_HOST>` | `your-idp.example.com` |

> ⚠️ **Firebase 注意**: Firebase の securetoken は公開鍵を **x509(PEM)** で配布し JWKS 非対応。
> Envoy `remote_jwks` は JWKS を要求するため **Firebase は直接使えない**。
> JWKS を公開する IdP(Auth0 / Keycloak / Cloudflare Access / Google Sign-In 等)を使うこと。
>
> ⚠️ **`sub` = ULID 前提**: `gateway.Parse` は `sub` を ULID として解釈する。
> IdP の `sub` が UUID 等の場合は sub→内部 ULID のマッピングが別途必要(現状未対応)。

## ローカル動作確認(任意)

ローカルは通常 Envoy を挟まず、`x-endpoint-api-userinfo` を手で埋め込んで app を直接叩く
(`../README.md` の gRPC API 参照)。Envoy 込みで試す場合は本物の IdP と JWT が要る。
