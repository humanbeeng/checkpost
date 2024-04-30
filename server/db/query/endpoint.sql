-- name: GetEndpoint :one
SELECT
    *
FROM
    "endpoint"
WHERE
    endpoint = $1
    AND is_deleted = FALSE
LIMIT
    1;

-- name: GetUserEndpoints :many
SELECT
    *
FROM
    "endpoint"
WHERE
    user_id = $1
    AND is_deleted = FALSE;

-- name: CheckEndpointExists :one
SELECT
    EXISTS (
        SELECT
            *
        FROM
            endpoint
        WHERE
            endpoint = $1
            AND expires_at > NOW()
            AND is_deleted = FALSE
        LIMIT
            1
    );

-- name: GetNonExpiredEndpointsOfUser :many
SELECT
    *
FROM
    "endpoint"
WHERE
    user_id = $1
    AND expires_at > NOW()
    AND is_deleted = FALSE;

-- name: InsertEndpoint :one
INSERT INTO
    endpoint (endpoint, user_id, plan, expires_at)
VALUES
    ($1, $2, $3, $4)
RETURNING
    *;

-- name: InsertFreeEndpoint :one
INSERT INTO
    endpoint (endpoint, user_id, plan, expires_at)
VALUES
    ($1, $2, 'free', $3)
RETURNING
    *;