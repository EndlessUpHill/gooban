module github.com/EndlessUpHill/gooban

go 1.22.5

// replace github.com/EndlessUpHill/goakka/core v0.0.0 => ../goakka/core

require (
	github.com/EndlessUpHill/goakka/core v0.0.4
	github.com/google/uuid v1.6.0
	github.com/jackc/pgx/v5 v5.7.1
)

require github.com/pkg/errors v0.9.1 // indirect

require (
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/pgx v3.6.2+incompatible
	github.com/jackc/puddle/v2 v2.2.2 // indirect
	golang.org/x/crypto v0.28.0 // indirect
	golang.org/x/sync v0.8.0 // indirect
	golang.org/x/text v0.19.0 // indirect
)
