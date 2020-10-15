#EMTCT

## Development Setup
Install docker compose.
Bring up the postgres dev database using docker compose:
`docker-compose up`

Run the migrations. This will create all the database tables:
`migrate -path db/migrations -database "postgres://postgres:password@localhost:5432/emtct?sslmode=disable" up`

## Front End
Start the front end in development mode: `NODE_ENV=development yarn start`.
This will read the environment variables from `.env.development`.
