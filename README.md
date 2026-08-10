# grpc-cqrs-go

Feature-first + CQRS アーキテクチャで構築した gRPC サーバーです。
レシピとユーザーの管理、および ReBAC による公開範囲・共有制御を行う API を提供します。

## 技術スタック

- **Go 1.25** - アプリケーション言語
- **gRPC** - API プロトコル
- **PostgreSQL 18** - データベース
- **Wire** - 依存性注入 (DI)
- **sqlc** - SQL からの型安全なコード生成
- **Atlas** - データベースマイグレーション
- **Buf** - Protobuf コード生成
- **ESPv2 (ESP) + Firebase Authentication** - 認証（前段のゲートウェイで検証しヘッダ伝播）
- **zerolog** - 構造化ログ

## アーキテクチャ

Feature-first + CQRS パターンを採用し、ドメイン（フィーチャ）ごとに Command（書き込み）と Query（読み取り）を分離しています。

### ディレクトリ構成

```
.
├── cmd/
│   └── serve/
│       ├── main.go                 # サーバーエントリーポイント
│       ├── wire.go                 # Wire DI 設定
│       └── wire_gen.go             # Wire 生成コード
├── database/
│   ├── migrations/                 # Atlas マイグレーションファイル
│   ├── queries/                    # sqlc クエリ定義 (recipe/tuple/user.sql)
│   └── schema.pg.hcl               # Atlas 宣言的スキーマ
├── internal/
│   ├── app/                        # アプリケーション層（フィーチャ単位）
│   │   ├── recipe/                 # Recipe ドメイン
│   │   │   ├── command/            #   書き込み (Create, Update, UpdateVisibility)
│   │   │   ├── query/              #   読み取り (Get)
│   │   │   ├── entity/             #   エンティティ (Recipe, Visibility)
│   │   │   ├── handler.go          #   gRPC ハンドラー
│   │   │   └── wire.go
│   │   ├── user/                   # User ドメイン
│   │   │   ├── command/            #   書き込み (CreateUser)
│   │   │   ├── entity/             #   エンティティ (User)
│   │   │   ├── handler.go
│   │   │   └── wire.go
│   │   ├── share/                  # Share 機能（レシピ共有）
│   │   │   ├── command.go          #   ShareRecipe（tuple 付与）
│   │   │   ├── handler.go
│   │   │   └── wire.go
│   │   ├── registrar.go            #   サービス登録
│   │   └── wire.go
│   ├── identity/                   # 現在ユーザーの identity（ctx 伝播・grpc 非依存）
│   │   ├── type.go                 #   UIDKey / UserID(ctx) / ErrUnauthenticated
│   │   └── gateway/                #   ゲートウェイヘッダの解釈 (Parse, HeaderUserInfo)
│   ├── authz/                      # 認可 (ReBAC)
│   │   ├── check.go                #   Checker (Can*/Is*)
│   │   ├── model.go                #   Tuple, ObjectType, Relation, Permission
│   │   ├── storage.go              #   Storage インターフェース
│   │   ├── error.go
│   │   └── wire.go
│   ├── db/                         # データベース層（エンティティ単位）
│   │   ├── recipe/                 #   recipe 永続化 (command + query)
│   │   ├── user/                   #   user 永続化
│   │   ├── tuple/                  #   relation_tuples 永続化
│   │   ├── gen/                    #   sqlc 生成コード (DO NOT EDIT)
│   │   ├── pool.go                 #   コネクションプール
│   │   └── wire.go
│   ├── grpcerr/                    # gRPC エラー変換 (WithStatus: cause 保持 + 安全な status)
│   ├── interceptor/                # gRPC インターセプター
│   │   ├── auth.go                 #   認証（ヘッダ→UID, GatewayAuthUnaryInterceptor + 公開メソッド除外）
│   │   ├── error.go                #   エラー処理（内部エラー秘匿 + cause ログ）
│   │   ├── logging.go              #   リクエストログ
│   │   └── recovery.go             #   パニックリカバリー
│   └── logger/
│       └── zerolog.go
├── espv2/                          # ESPv2 (API ゲートウェイ) 導入 scaffold (config/compose/README)
├── pkg/pb/                         # Protobuf 生成コード (recipe/share/user)
├── proto/                          # Protobuf 定義 (recipe/share/user)
├── atlas.hcl                       # Atlas 設定
├── buf.yaml / buf.gen.yaml         # Buf 設定
├── compose.yaml                    # Docker Compose
├── mise.toml                       # ツールバージョン管理
└── sqlc.yaml                       # sqlc 設定
```

### CQRS パターン

各ドメイン（フィーチャ）は以下のレイヤーで構成されています:

```
handler.go          ← gRPC リクエストの受付・バリデーション・エラー変換
  ├── command/      ← 書き込み操作（エンティティ操作 → Storage インターフェース）
  ├── query/        ← 読み取り操作（Storage → 読み取りモデル）
  └── entity/       ← エンティティの定義（Recipe / User など）
```

- Storage インターフェースにより app 層と DB 層が疎結合（依存の向きは `db → app`）
- DB 層はエンティティ単位（`db/recipe`, `db/user`, `db/tuple`）で実装。CQRS の read/write の差は app 側の interface が表現する
- エンティティは非公開フィールド + コンストラクタ + getter でカプセル化（不正な生成を型で防止）

### 認証 / identity

app 内でトークン検証はしない。**前段(ゲートウェイ)が検証済みの ID をヘッダで渡す前提**で、
app はそれを信頼して現在ユーザーの identity を確定・伝播するだけ（`internal/identity`）。

- **本番**: ESP(ESPv2/Envoy) が Firebase の ID トークンを検証し、`X-Endpoint-API-UserInfo`
  ヘッダを付与 → `GatewayAuthUnaryInterceptor` が UID を取り出し ctx へ
- **ローカル**: ESP 無し。開発者がヘッダを直接埋め込む（検証なし・impersonate 用。`espv2/README.md` 参照）
- 以降 `identity.UserID(ctx)` で `ulid.ULID` を取得（ctx を読むのは handler だけ、下層へは typed で渡す）
- **前提**: Firebase の uid がアプリの user id（ULID）であること。`CreateUser` は認証済み UID を user.id として登録する
- reflection / health は認証不要（`GatewayAuthUnaryInterceptor` が除外。health は本番の probe 用）

> ⚠️ app はヘッダを信頼するため、**ゲートウェイ経由でのみ到達可能**にすること
> （直接到達できるとヘッダ偽造で認証を素通りできる）。

### 認可 (ReBAC)

Relationship-Based Access Control を採用。

- **Visibility**: `public`（誰でも） / `private`（owner のみ） / `restricted`（owner + 共有相手）
- **Relation**: `owner ⊃ editor ⊃ viewer`（owner は複数可）
- **Tuple ストア**: `relation_tuples` テーブルで「誰が・何に・どの relation か」を管理。subject は不透明な文字列として扱い、ULID かどうかは authz の関心外
- **Permission**: `map[ObjectType]map[Permission][]Relation` で action → 許可 relation を定義
- view の認可ゲートは app 層で分岐、edit/share/公開範囲変更は owner/editor を tuple で判定

### エラー処理

- ハンドラは既知エラーを `grpcerr.WithStatus(cause, code, msg)` で変換（client には安全な status、ログには cause を保持）
- 未変換の生エラーは `ErrorHandlingUnaryInterceptor` が汎用 `Internal` + cause ログに（**内部エラーを client に漏らさない**）
- 横断的な `PermissionDenied` は interceptor で集約
- DB 境界で低レベルエラーを翻訳（`23505` → AlreadyExists / `pgx.ErrNoRows` → NotFound）

## セットアップ

### 前提条件

- [mise](https://mise.jdx.dev/) (ツールバージョン管理)
- Docker / Docker Compose

### 1. ツールのインストール

```bash
mise install
```

### 2. データベースの起動

```bash
docker compose up -d
```

PostgreSQL が `localhost:25432` で起動します。

### 3. マイグレーションの実行

```bash
atlas migrate apply --env local
```

### 4. サーバーの起動

```bash
go run ./cmd/serve
```

サーバーが `localhost:50051` で起動します。ポートは `-port` フラグで変更可能です。

## gRPC API

認証は前段のゲートウェイが検証済み ID を `x-endpoint-api-userinfo` ヘッダで渡す前提（[認証 / identity](#認証--identity) 参照）。

- **本番**（ESP 前段）: クライアントは `authorization: Bearer <firebase_id_token>` を送る（ESP が検証しヘッダ付与）
- **ローカル**（ESP 無し）: `x-endpoint-api-userinfo` を手で埋め込む。以下の例はこの形。

```bash
# 対象ユーザーの ULID を base64url(JSON) にしてヘッダ値に（app は padding 有無どちらも可）
UID=01J...
INFO=$(printf '{"sub":"%s"}' "$UID" | basenc --base64url | tr -d '=')
```

### UserService

#### CreateUser（要認証）

認証済みの Firebase ユーザーを app 側の users に登録します（user id = Firebase uid）。

```bash
grpcurl -plaintext -H "x-endpoint-api-userinfo: $INFO" -d '{
  "name": "Kodai",
  "email": "kodai@example.com"
}' localhost:50051 user.UserService/CreateUser
```

### RecipeService

#### CreateRecipe（要認証）

作成者が owner として登録されます（`private` 初期）。

```bash
grpcurl -plaintext -H "x-endpoint-api-userinfo: $INFO" -d '{
  "title": "カレーライス",
  "description": "スパイスから作る本格カレー"
}' localhost:50051 recipe.RecipeService/CreateRecipe
```

#### UpdateRecipe（要認証・editor 以上）

```bash
grpcurl -plaintext -H "x-endpoint-api-userinfo: $INFO" -d '{
  "id": "<recipe_id>",
  "title": "カレーライス 改",
  "description": "改良版"
}' localhost:50051 recipe.RecipeService/UpdateRecipe
```

#### ChangeVisibility（要認証・owner のみ）

```bash
grpcurl -plaintext -H "x-endpoint-api-userinfo: $INFO" -d '{
  "id": "<recipe_id>",
  "visibility": "VISIBILITY_PUBLIC"
}' localhost:50051 recipe.RecipeService/ChangeVisibility
```

#### GetRecipe

visibility に応じて認可されます（public は誰でも / private は owner / restricted は共有相手）。

```bash
grpcurl -plaintext -H "x-endpoint-api-userinfo: $INFO" -d '{
  "id": "<recipe_id>"
}' localhost:50051 recipe.RecipeService/GetRecipe
```

レスポンス:
```json
{
  "recipe": {
    "id": "550e8400-...",
    "userId": "01JNQF...",
    "title": "カレーライス",
    "description": "スパイスから作る本格カレー",
    "visibility": "VISIBILITY_PRIVATE",
    "createdAt": "2026-03-08T...",
    "updatedAt": "2026-03-08T..."
  },
  "user": { "id": "01JNQF...", "name": "Kodai", "email": "kodai@example.com" }
}
```

### ShareService

#### ShareRecipe（要認証・owner のみ）

対象ユーザーに `viewer` / `editor` を付与します。重複付与は `AlreadyExists` になります。

```bash
grpcurl -plaintext -H "x-endpoint-api-userinfo: $INFO" -d '{
  "recipe_id": "<recipe_id>",
  "target_user_id": "<user_id>",
  "relation": "viewer"
}' localhost:50051 share.ShareService/ShareRecipe
```

### サービス一覧の確認

gRPC リフレクションが有効です:

```bash
grpcurl -plaintext localhost:50051 list
```

## コード生成

```bash
# Protobuf → Go コード生成
buf generate

# SQL → Go コード生成
sqlc generate

# Wire DI コード生成
cd cmd/serve && wire

# スキーマ変更からマイグレーション生成 (Atlas)
atlas migrate diff <name> --env local
```
