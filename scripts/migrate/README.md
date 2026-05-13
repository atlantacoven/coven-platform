# migrate

This Go CLI manages database migrations for the member-site backend. Migrations are stored in `member-site/migrations/`.

```bash
# create new
go run ./scripts/migrate create "add_my_table"
# run all pending migrations
go run ./scripts/migrate up
```
