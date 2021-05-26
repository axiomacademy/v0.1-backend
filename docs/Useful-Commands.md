# Look up lifesaving commands :tools: 

1) `go get -u ./...`
    * This command lets you pull all your essential go packages in one fell swoop, essential after a fresh clone
2) `go run github.com/99designs/gqlgen generate`
    * Runs the gqlgen generate code, such that the resolves are automatically generated
3) `gofmt -s -w .`
    * Runs go format on the code, you should be doing this as part of pre-commit hooks, on save or you should at least be running it before you commit
4) `migrate -source file://path/to/migrations -database postgres://localhost:5432/database up 2`
    * Migrate to a specific database version (you can replace up with down)
5) `migrate -source file://path/to/migrations -database postgres://localhost:5432/database drop`
    * Drop entire database, useful when you reach a dirty database state