
# BLOGESTON
A simple backend blog implemented with Golang and Postgres

## Environment Variables

To run this project, you will need to add the following environment variables to your .env file

`DB_CONNECTION`

`JWT_SECRET_KEY`
## Run & Build App

Run this command in the terminal when you are in the project directory `/blogeston`.

```bash
  go run main.go
```

Using the command below in the project directory `/blogeston`, you can create an executable binary for our sample Go application, which allows you to distribute and run the application wherever you want.
```bash
  go build
```
## API Reference

### Public API

#### Register user
```http
  POST /register
```
#### Login user
```http
  POST /login
```


#### Get all Posts

```http
  GET /posts
```


#### Get a Post

```http
  GET /posts/:id
```

| Parameter | Type     | Description                |
| :-------- | :------- | :------------------------- |
| `id` | `int` | **Required** Id of item to fetch |

#### Get comments on a post

```http
  GET /posts/:id/comments
```
| Parameter | Type     | Description                |
| :-------- | :------- | :------------------------- |
| `id` | `int` | **Required** Id of item to fetch |


### Private API

#### Create post

```http
  POST /create-post
```

| Parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `authorization`      | `string` | **Required**. to authorize user |


#### Create comment on a post

```http
  POST /posts/:id/create-comment
```

| Parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `id`      | `int` | **Required**. Id of item to fetch |
| `authorization`      | `string` | **Required**. to authorize user |


#### React to comment on a post

```http
  POST /comments/:comment_id/react
```

| Parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `comment_id`      | `int` | **Required**. Id of item to fetch |
| `authorization`      | `string` | **Required**. to authorize user |


### Private API for implement Role Based Access Control

#### Get All Users for admin role
```http
  GET /users
```

| Parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `authorization`      | `string` | **Required**. to authorize user |

#### Publish post
```http
  POST /posts/:id/publish
```

| Parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `id`      | `int` | **Required**. Id of item to fetch |
| `authorization`      | `string` | **Required**. to authorize user |

## Acknowledgements
 - [Go programming language](https://go.dev/)
 - [Gin - Go Web Framework](https://github.com/gin-gonic/gin)
 - [PostgreSQL](https://www.postgresql.org/)


