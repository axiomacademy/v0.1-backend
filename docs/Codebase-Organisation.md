# Codebase Organisation 📁
This section will tell you some ways in which the codebase has been organised. In many ways, we use a lot of intuitive rules and personal biases to organise our codebases. This hopes to clarify those biases and put them in words so that it is understandable for future developers. It will go through the file structure and some rules of thumb too.

## File Structure
The file tree of the backend repository currently looks like this
```bash
.
├── README.md
├── db
│   ├── affinity.go
│   ├── availability.go
│   ├── database.go
│   ├── lesson.go
│   ├── match.go
│   ├── migrations
│   │   ├── 000001_create_students_table.down.sql
│   │   ├── 000001_create_students_table.up.sql
│   │   ├── 000002_create_tutors_table.down.sql
│   │   ├── 000002_create_tutors_table.up.sql
│   │   ├── 000003_create_lessons_table.down.sql
│   │   ├── 000003_create_lessons_table.up.sql
│   │   ├── 000004_create_affinity_table.down.sql
│   │   ├── 000004_create_affinity_table.up.sql
│   │   ├── 000005_create_notifications_table.down.sql
│   │   └── 000005_create_notifications_table.up.sql
│   ├── notification.go
│   ├── student.go
│   ├── subject.go
│   ├── tutor.go
│   └── utilities.go
├── go.mod
├── go.sum
├── gqlgen.yml
├── graph
│   ├── generated
│   │   └── generated.go
│   ├── model
│   │   └── models_gen.go
│   ├── resolver.go
│   ├── schema.graphqls
│   ├── schema.resolvers.go
│   └── utilities.go
├── middlewares
│   ├── auth.go
│   └── generic.go
├── server.go
├── services
│   ├── match
│   │   └── match.go
│   └── notifs
│       └── notifs.go
└── utilities
    └── auth
        └── auth.go
```

Of course, this tree will only keep on growing so I will attempt to go through some of the basic features of this filesystem organisation method.

## `server.go` and other top level files :top: 
These are the files you see in the root directory. Files like `server.go` and `gqlgen.yml` are all core program files. `server.go` as you might have guessed contains the entry point for our backend code, it sets up all the essential services, runs the main logging instance and runs the main HTTP handler. The other files are configuration files intended to modify the behaviour of Go itself or Go libraries.

## `db` :book: 
This is the folder which contains all the database relevant files. This is referred to as the `Repository` in internal code variable and contains all the SQL queries to be found in the backend. In fact, the Postgres DB handler instance never leaves the Repository, and **should never** leave the `Repository`. 

Inside, it contains relevant SQL queries wrapped in readable functions with appropriate parameters. A general rule of thumb is to limit each function to exactly **one** SQL query and if there are multiple queries they should be split up into multiple functions. Another rule is that, repository functions should never output logging statements. They should simply propagate their errors upwards to be handled by the calling functions.

The `db` folder also contains migrations which are handled by [golang-migrate](https://github.com/golang-migrate/migrate). You can look at [:tractor: Getting Started](getting-started) for more information on applying and creating migrations.

## `graph` :chart_with_upwards_trend: 
This folder contains most of the GraphQL resolvers. Some important things to note would be that `generated.go` and `models_gen.go` are generated files and should never be manually touched. The only files that are of concern are any `.graphqls` files and their associated `.resolvers.go` files. 

GraphQL resolvers are in charge of calling the appropriate repository and service methods. The also should solely determine any error message returned by the GraphQL API.

The `.graphqls` files contain the GraphQL schemas, written in the GraphQL SDL. The `.resolvers.go` files contain the resolver implementation for those schemas and are partially generated and managed by gqlgen. You can read more about it on their [github page](https://github.com/99designs/gqlgen).

## `middlewares`
This contains HTTP middleware implementations. The only one of interest is `auth.go` which implements the authentication middleware, taking the JWT in the `token` cookie and then parsing it to a `Student` or `Tutor` type.

## `services` :robot: 
Services are what tie your database layer to your resolvers. In instances where a resolver has to perform business logic more complicated than a simple SQL request, a service should be written to abstract out that logic. Services are allowed and are encouraged to take a logger instance and use it to log any errors that might be generated. They should also pass their errors up to their calling functions which are typically the main GraphQL resolvers.

## `utilities` :tools: 
Utilities are just general-purpose functions, used widely throughout the program. They have been extracted and moved outside to their own package to maintain DRY principles.