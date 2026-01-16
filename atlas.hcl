data "external_schema" "gorm" {
  program = [
    "go",
    "run",
    "./packages/database/cmd/loader",
  ]
}

locals {
  db_url = "postgresql://${getenv("DB_USER")}:${getenv("DB_PASSWORD")}@${getenv("DB_HOST")}:${getenv("DB_PORT")}/${getenv("DB_NAME")}?sslmode=disable"
}

env "gorm" {
  src = data.external_schema.gorm.url
  dev = local.db_url
  url = local.db_url
  migration {
    dir = "file://packages/database/migrations"
  }
  format {
    migrate {
      diff = "{{ sql . \"  \" }}"
    }
  }
}