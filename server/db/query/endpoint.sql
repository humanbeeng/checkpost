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

-- name: InsertEndpoint :one
insert into
    endpoint (
        endpoint, user_id, plan, expires_at
    )
values ($1, $2, $3, $4)
returning
    *;

-- name: InsertGuestEndpoint :one
insert into
    endpoint (endpoint, expires_at, plan)
values ($1, $2, 'guest')
returning
    *;

-- name: InsertFreeEndpoint :one
insert into
    endpoint (endpoint, user_id, plan, expires_at)
values ($1, $2, 'free', $3)
returning
    *;