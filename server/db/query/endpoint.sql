-- name: GetEndpoint :one
select *
from "endpoint"
where
    endpoint = $1
    and is_deleted = false
limit 1;

-- name: CheckEndpointExists :one
select exists (
        select *
        from endpoint
        where
            endpoint = $1
            and expires_at > now()
            and is_deleted = false
        limit 1
    );

-- name: GetNonExpiredEndpointsOfUser :many
select *
from "endpoint"
where
    user_id = $1
    and expires_at > now()
    and is_deleted = false;

-- name: CreateNewEndpoint :one
insert into
    endpoint (
        endpoint, user_id, plan, expires_at
    )
values ($1, $2, $3, $4)
returning
    *;

-- name: CreateNewGuestEndpoint :one
insert into
    endpoint (endpoint, expires_at)
values ($1, $2)
returning
    *;

-- name: CreateNewFreeUrl :one
insert into
    endpoint (endpoint, user_id, expires_at)
values ($1, $2, $3)
returning
    *;