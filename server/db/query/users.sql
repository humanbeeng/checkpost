-- name: CreateUser :one
insert into users (
	name, plan ,email
) values ($1, $2, $3) returning *;


-- name: GetUser :one
select * from users where id = $1 limit 1; 

-- name: ListUsers :many
select * from users limit $1 offset $2;

-- name: DeleteUser :exec
delete from users where id = $1;
