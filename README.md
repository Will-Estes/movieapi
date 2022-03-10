# movieapi
Simple movie api built with Go. It uses Postgres as its database and is capable of simple JWT token functions.

# Env Variables
You will need an envirnment variable for JWT_SECRET and DB_CONNECTION (both found in main.go). 
The JWT_SECRET is for decoding your JSON Web Tokens. The DB_CONNECTIOn is for making a connection to a Postgres database

# Binary for Linus (Ubuntu)
`env GOOS=linux GOARCH=amd64 go build -o gomovies ./cmd/api`
