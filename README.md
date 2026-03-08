# grpc-cqrs-go

Feature-first + CQRS アーキテクチャで構築した gRPC サーバーです。
レシピとユーザーの管理を行う API を提供します。

## 技術スタック

- **Go 1.25** - アプリケーション言語
- **gRPC** - API プロトコル
- **PostgreSQL 18** - データベース
- **Wire** - 依存性注入 (DI)
- **sqlc** - SQL からの型安全なコード生成
- **Atlas** - データベースマイグレーション
- **Buf** - Protobuf コード生成
- **zerolog** - 構造化ログ

## アーキテクチャ

Feature-first + CQRS パターンを採用し、ドメインごとに Command（書き込み）と Query（読み取り）を分離しています。

### ディレクトリ構成

```
.
├── cmd/
│   └── serve/
│       ├── main.go                 # サーバーエントリーポイント
│       ├── wire.go                 # Wire DI 設定
│       └── wire_gen.go             # Wire 生成コード
├── database/
│   ├── initdb.d/
│   │   └── init.sql                # DB 初期化スクリプト
│   ├── migrations/                 # Atlas マイグレーションファイル
│   └── queries/                    # sqlc クエリ定義
│       ├── recipe.sql
│       └── user.sql
├── internal/
│   ├── app/                        # アプリケーション層
│   │   ├── recipe/                 # Recipe ドメイン
│   │   │   ├── command/            #   書き込み操作 (CreateRecipe)
│   │   │   ├── query/              #   読み取り操作 (GetRecipe)
│   │   │   ├── domain/             #   ドメインエンティティ
│   │   │   ├── handler.go          #   gRPC ハンドラー
│   │   │   └── wire.go             #   DI 設定
│   │   ├── user/                   # User ドメイン
│   │   │   ├── command/            #   書き込み操作 (CreateUser)
│   │   │   ├── domain/             #   ドメインエンティティ
│   │   │   ├── handler.go          #   gRPC ハンドラー
│   │   │   └── wire.go             #   DI 設定
│   │   ├── registrar.go            #   サービス登録
│   │   └── wire.go
│   ├── db/                         # データベース層
│   │   ├── command/                #   書き込み用ストレージ実装
│   │   ├── query/                  #   読み取り用ストレージ実装
│   │   ├── gen/                    #   sqlc 生成コード
│   │   ├── pool.go                 #   コネクションプール
│   │   └── wire.go
│   ├── interceptor/                # gRPC インターセプター
│   │   ├── logging.go              #   リクエストログ
│   │   └── recovery.go             #   パニックリカバリー
│   └── logger/
│       └── zerolog.go              #   ロガー初期化
├── pkg/pb/                         # Protobuf 生成コード
│   ├── recipe/
│   └── user/
├── proto/                          # Protobuf 定義
│   ├── recipe/
│   │   └── recipe.proto
│   └── user/
│       └── user.proto
├── atlas.hcl                       # Atlas 設定
├── buf.yaml                        # Buf 設定
├── buf.gen.yaml                    # Buf コード生成設定
├── compose.yaml                    # Docker Compose
├── mise.toml                       # ツールバージョン管理
└── sqlc.yaml                       # sqlc 設定
```

### CQRS パターン

各ドメインは以下のレイヤーで構成されています:

```
handler.go          ← gRPC リクエストの受付・バリデーション
  ├── command/      ← 書き込み操作（ドメインロジック → Storage インターフェース）
  ├── query/        ← 読み取り操作（Storage インターフェース → ドメインモデル）
  └── domain/       ← ドメインエンティティの定義
```

Storage インターフェースにより、ドメイン層とデータベース層が疎結合になっています。

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

```bash
go run ./cmd/serve -port 8080
```

## gRPC API

### UserService

#### CreateUser

```bash
grpcurl -plaintext -d '{
  "name": "Kodai",
  "email": "kodai@example.com"
}' localhost:50051 user.UserService/CreateUser
```

レスポンス:
```json
{
  "userId": "01JNQF..."
}
```

### RecipeService

#### CreateRecipe

```bash
grpcurl -plaintext -d '{
  "userId": "<user_id>",
  "title": "カレーライス",
  "description": "スパイスから作る本格カレー"
}' localhost:50051 recipe.RecipeService/CreateRecipe
```

レスポンス:
```json
{
  "recipeId": "550e8400-..."
}
```

#### GetRecipe

```bash
grpcurl -plaintext -d '{
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
    "createdAt": "2026-03-08T...",
    "updatedAt": "2026-03-08T..."
  },
  "user": {
    "id": "01JNQF...",
    "name": "Kodai",
    "email": "kodai@example.com"
  }
}
```

### サービス一覧の確認

gRPC リフレクションが有効なため、以下のコマンドでサービス一覧を確認できます:

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
```
