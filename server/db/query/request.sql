-- name: CreateNewRequest :one
insert into
  request (
    user_id,
    endpoint_id,
    path,
    response_id,
    content,
    method,
    source_ip,
    content_size,
    response_code,
    headers,
    query_params
  )
values
  ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) returning *;


-- name: GetEndpointHistory :many
select * from request where endpoint_id = $1 limit $2 offset $3;


-- name: GetRequestById :one
select * from request where id = $1 limit 1;
