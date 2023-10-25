table "users" {
  schema = schema.jab

  column "id" {
    type = serial
    null = false
  }

  column "network_id" {
    type = varchar(255)
    null = false
  }

  column "name" {
    type = varchar(255)
    null = false
  }

  column "permissions" {
    type = sql("varchar[]")
    null = false
  }

  primary_key {
    columns = [
      column.id
    ]
  }

  index "users.network_id.idx" {
    columns = [
      column.network_id
    ]
    unique = true
  }
}
