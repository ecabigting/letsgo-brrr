# letsgo-brrr
A collection of APIs, Microservice and Web Apps in general using Go


## ℹ️  devicemonitor
A web app build with `go`, `websockets`, and `htmx` to show realtime device info. 
Info is read using the using the [gopsutil](https://github.com/shirou/gopsutil)

[devicemonitor](/devicemonitor/src.png)

### Installation 
Clone this repo and `cd` into the project directory
```zsh
$ git clone https://github.com/ecabigting/letsgo-brrr.git
$ cd kanbanapi
```
Run the app by executing
```zsh
$ go run ./cmd
```

Open a browser and visit the URL: `http://localhosts:<port>`


> The UI is build with bootstrap


## 📝  kanbanapi
A simple rest `api` with [xata](https://xata.io/) a `postgresql` database

[kanbanapi](/kanbanapi/src.png)

### Installation 
Clone this repo and `cd` into the project directory
```zsh
$ git clone https://github.com/ecabigting/letsgo-brrr.git
$ cd kanbanapi
```
Create your `.env` file in the root directory `kanbanapi/` and add the following:
```zsh
#.env
XATA_PSQL_URL=
JWT_TOKEN_SECRET=
```
 - Set the your `postgresql` url to `XATA_PSQL_URL` in my case I am using `xata` as my database
 - Set your chosen secret to `JWT_TOKEN_SECRET`

Run the project with:
```zsh
go run ./
```
Alternatively you can use `air` for hot reload module. You can install `air` by following these [instructions](https://github.com/air-verse/air?tab=readme-ov-file#installation)
#### Available endpoints
You can interact with the API using your favorite endpoints by using the following URLS:

- `/` the root URL
- `/register` a `POST` request to register a user and returns an access token
- `/login` a `POST` to login a registered user and returns an access token
- `/projects` a `GET` request returns all projects owned by the requesting user, required an access token in the `Authorization` header
- `/projects/{xata_id}` a `GET` request that returns the project the of the given `xata_id`, project must be owned by the requesting user, required an access token in the `Authorization` header
- `/projects/{xata_id}` a `PUT` request to update the project related to the given `xata_id`, required an access token in the `Authorization` header
- `/projects` a `POST` request to create a new project, required an access token in the `Authorization` header
at- `/projects/{xata_id}` a `DELETE` request to delete a project related to the given `xata_id`, required an access token in the `Authorization` header

> Check the `/kanbanapi/request-endpoints` for the collection of request endpoints with body 





