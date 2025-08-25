resource "random_password" "db" {
  length  = 24
  special = false
}

resource "aws_db_subnet_group" "rds" {
  name       = "${local.name}-db-subnets"
  subnet_ids = [for s in aws_subnet.private : s.id]
}

resource "aws_db_instance" "pg" {
  identifier              = "${local.name}-pg"
  backup_retention_period = 0
  engine                  = "postgres"
  engine_version          = "16"
  instance_class          = "db.t4g.micro"
  allocated_storage       = 20
  db_name                 = "appdb"
  username                = "app"
  password                = random_password.db.result
  db_subnet_group_name    = aws_db_subnet_group.rds.name
  vpc_security_group_ids  = [aws_security_group.db.id]
  publicly_accessible     = false
  skip_final_snapshot     = true
  deletion_protection     = false
}
