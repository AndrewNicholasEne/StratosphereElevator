resource "aws_secretsmanager_secret" "db" {
  name = "${local.name}/db"
}

resource "aws_secretsmanager_secret_version" "db" {
  secret_id = aws_secretsmanager_secret.db.id
  secret_string = jsonencode({
    username = aws_db_instance.pg.username
    password = random_password.db.result
    host     = aws_db_instance.pg.address
    port     = 5432
    database = aws_db_instance.pg.db_name
    url      = "postgres://${aws_db_instance.pg.username}:${random_password.db.result}@${aws_db_instance.pg.address}:5432/${aws_db_instance.pg.db_name}?sslmode=require"
  })
}
