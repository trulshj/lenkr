schema "main" {}

table "links" {
  schema = schema.main
  column "url" {
    null = false
    type = text
  }
  column "short_url" {
    null = false
    type = text
  }
  column "fetches" {
    null = false
    type = integer
    default = 0
  }
  primary_key {
    columns = [column.short_url]
  }
}
