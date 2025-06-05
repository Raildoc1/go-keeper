# Password Manager "Go Keeper"
## Server
### Enviroment Variables
- `DATABASE_URI` (Required) — PostgreSQL database connection string \
  (example `postgres://user:password@localhost:5432/gokeeper?sslmode=disable`)
- `SERVER_ADDRESS` (Optional, default is `":8080"`) — the address at which the server will be running
## Client
### Enviroment Variables
- `SERVER_ADDRESS` (Optional, default is `"local:8080"`) — the address of a server to connect to
### Usage
#### Authentication

- type `register` to start new account creation process or `login` if you already have one and hit enter then enter login and password
- you will stay logged in until received token is expired

#### Storage Management

type one of the following commands and follow instructions printed:

- `list` prints all locally stored entries' metadata and specifies whether they are stored on the server
- `load` prints chosen entry content
- `store` starts new entry creation process **(this command puts entry in local storage only, you should then call `sync` command)**
- `sync` synchronizes local storage with the server's one
- `logout` clears token and returns to Authentication stage
- `quit` quits the application

## Docker

```
docker compose build
docker compose run --rm client
```