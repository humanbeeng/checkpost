-- name: CreateNewEndpoint :one
insert into
  endpoint (endpoint, user_id, plan, expires_at)
values
($1, $2, $3, $4) returning *;

-- name: CreateNewGuestEndpoint :one
insert into endpoint (endpoint, expires_at) values ($1, $2) returning *;

-- name: CheckEndpointExists :one
select exists(select * from endpoint where endpoint = $1 limit 1);

-- name: GetEndpoint :one
select * from "endpoint" where endpoint = $1 limit 1;
