# blog-aggregator
RSS blog aggregator service in Go

You'll probably want to create a .env file with these environment variables:
- PORT=port (e.g. 8080)
- DB_CON_STRING=protocol://username:password@host:port/database?sslmode=disable (e.g. postgres://postgres:12345@localhost:5432/blogs?sslmode=disable)