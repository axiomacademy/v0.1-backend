# Axiom Backend - Quick Start Guide
First of all, welcome to the team :smile: let's get you up and running in a jiffy! To get the core backend code running on your laptop is fairly straightforward, all you need to be aware of are some of the basic tools we use in the Axiom backend. 

## Our core tools :tools: 
1) **Git** - our favourite version control, make sure you pick that up (paired with Gitlab)
2) **Golang** - head on over [here](https://golang.org/) if you need help getting that installed
3) **Go modules** - makes go so much easier to deal with
4) **gqlgen** - our [GraphQL](https://graphql.org/) client, uses go generate to create graphql resolvers as you will see in the future steps
5) **Postgres** - the SQL database of choice
6) **Docker** - containerise all the things :whale:, grab that [here](https://www.docker.com/)
7) **docker-compose** - we use this for rudimentary container orchestration at the developer level

If you have trouble or don't understand any of the above mentioned tools entirely, feel free to take some time now to get used to them. Also approach any of our more senior developers, and they would be glad to lend you a hand :stuck_out_tongue_winking_eye: 

## The holy grail commands
Make sure once again that all the aformentioned tools are installed on your system, then follow these steps to start your own local development instance of the Axiom backend.

1) Run your own postgres instance using Docker :whale: 
```bash
docker run --name axiom -e POSTGRES_PASSWORD=password -d postgres
```
This creates the default postgres database, which is all we need to play around with. Specifically, it creates a database with the connection url `postgresql://postgres:axiom@127.0.0.1:5432/postgres?sslmode=disable`.

2) Populate your database with actual test data
```bash
psql run testData.sql
```
Having an empty database is kind of useless so let's populate it with some dummy data we can use later on when we call the GraphQL endpoints :wink: 

3) Set your environmental variables
* `PORT`: API access port, defaults to 8080
* `DB_URL`: Database connection URL, defaults to the connection URL shown aboce
* `SERVER_SECRET`: The server secret used to encode the JWT tokens, defaults to password
* `GOOGLE_APPLICATION_CREDENTIALS`: Firebase credentials for notifications

God yes, we know you hate these, we forget to set them all the time too :angry:. So it might be wise to add a script that sets all of these in one fell sweep, or to add it to your `.bashrc` or `.zshrc`. Be careful not to commit this script to the repository though, store it outside the repository directory!

4) Run the holy grail commands :thumbsup: 
```bash
git clone ssh://git@gitlab.solderneer.me:8002/axiom/backend.git
go get -u ./...

go run github.com/99designs/gqlgen generate
go run ./server.go
```
This clones and runs the instance of the backend which is on the `master` branch right now. Do note that the go generate call is what triggers the regeneration of the resolver stubs for GraphQL, so remember to run the generate call whenever you touch the GraphQL schema implementation, and as general practice whenever you push code to the repository.
