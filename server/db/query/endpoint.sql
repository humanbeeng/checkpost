-- name: CreateNewEndpoint :one
insert into
  endpoint (endpoint, user_id, plan)
values
($1, $2, $3) returning *;

-- name: CreateNewGuestEndpoint :one
insert into endpoint (endpoint) values ($1) returning *;

-- name: CheckEndpointExists :one
select exists(select * from endpoint where endpoint = $1 limit 1);

-- name: GetEndpoint :one
select * from "endpoint" where endpoint = $1 limit 1;
