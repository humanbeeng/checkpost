-- name: CreateNewRequest :one
INSERT INTO
    request (
        user_id,
        endpoint_id,
        PATH,
        response_id,
        CONTENT,
        METHOD,
        UUID,
        source_ip,
        content_size,
        response_code,
        headers,
        query_params,
        expires_at
    )
VALUES
    (
        $1,
        $2,
        $3,
        $4,
        $5,
        $6,
        $7,
        $8,
        $9,
        $10,
        $11,
        $12,
        $13
    )
RETURNING
    *;

-- name: GetEndpointHistory :many
SELECT
    request.id,
    request.uuid,
    request.user_id,
    request.plan,
    request.path,
    request.response_id,
    request.response_code,
    request.content,
    request.method,
    request.source_ip,
    request.content_size,
    request.headers,
    request.query_params,
    request.created_at,
    request.expires_at,
    endpoint.endpoint AS endpoint
FROM
    request
    LEFT JOIN endpoint ON request.endpoint_id = endpoint.id
WHERE
    endpoint.endpoint = $1
    AND request.user_id = $2
    AND request.is_deleted = FALSE
    AND request.expires_at > NOW()
ORDER BY
    request.id DESC
LIMIT
    $3
OFFSET
    $4;

-- name: GetRequestById :one
SELECT
    *
FROM
    request
WHERE
    id = $1
    AND is_deleted = FALSE
    AND expires_at > NOW()
LIMIT
    1;

-- name: GetEndpointRequestCount :one
SELECT
    COUNT(*) AS total_count,
    COUNT(
        CASE
            WHEN response_code = 200 THEN 1
        END
    ) AS success_count,
    COUNT(
        CASE
            WHEN response_code != 200 THEN 1
        END
    ) AS failure_count
FROM
    request r
    LEFT JOIN endpoint e ON r.endpoint_id = e.id
WHERE
    endpoint = $1
    AND is_deleted = FALSE
    AND expires_at > NOW();

-- name: GetRequestByUUID :one
SELECT
    *
FROM
    request
WHERE
    UUID = $1
    AND is_deleted = FALSE
    AND expires_at > NOW()
LIMIT
    1;

-- name: DeleteExpiredRequests :exec
DELETE FROM request
WHERE
    expires_at < NOW();