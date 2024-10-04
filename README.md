# letsgo-brrr
A collection of APIs, Microservice and Web Apps in general using Go

## 📝  kanbanapi
A simple rest `api` with the [xata](https://xata.io/) a `postgresql` database

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

