schema "gcook" {}

table "recipes" {
  schema = schema.gcook
  column "id" {
    type = uuid
  }
  column "user_id" {
    type = varchar(26)
  }
  column "title" {
    type = varchar(255)
  }
  column "description" {
    type = varchar(255)
  }
  column "created_at" {
    type = timestamp
  }
  column "updated_at" {
    type = timestamp
  }
  primary_key {
    columns = [column.id]
  }
}

