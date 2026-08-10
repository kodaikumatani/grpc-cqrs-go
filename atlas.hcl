env "local" {
  src = "file://./database/schema.pg.hcl"
  url = "postgres://gcook:p4ssw0rd!@localhost:25432/gcook?search_path=gcook&sslmode=disable"
  dev = "docker://postgres/18/dev?search_path=public"

  migration {
    dir = "file://./database/migrations"
  }

  format {
    migrate {
      diff = "{{ sql . \"  \" }}"
    }
  }
}

# コンテナ(k8s の migrate initContainer)からの `atlas migrate apply --env k8s` 用。
# 接続先は環境変数 DATABASE_URL から取る(シェル不要・secret を args に出さない)。
env "k8s" {
  url = getenv("DATABASE_URL")

  migration {
    dir = "file://./database/migrations"
  }
}
