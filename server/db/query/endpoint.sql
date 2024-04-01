-- name: CreateNewEndpoint :one
insert into
  endpoint (endpoint, user_id, plan)
values
($1, $2, $3) returning *;


-- name: GetEndpoint :one
select endpoint from "endpoint" where endpoint = $1 limit 1;
