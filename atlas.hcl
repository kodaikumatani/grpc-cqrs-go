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
