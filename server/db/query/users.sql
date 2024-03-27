-- name: CreateUser :one
insert into "user" (
	name, username, plan ,email
) values ($1, $2, $3, $4) returning *;


-- name: GetUser :one
select * from "user" where id = $1 limit 1; 

-- name: GetUserFromEmail :one
select * from "user" where email = $1 limit 1;

-- name: GetUserFromUsername :one
select * from "user" where username = $1 limit 1;

-- name: ListUsers :many
select * from "user" limit $1 offset $2;

-- name: DeleteUser :exec
delete from "user" where id = $1;
