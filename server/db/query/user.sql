-- name: CreateUser :one
INSERT INTO
	"user" (NAME, avatar_url, username, plan, email)
VALUES
	($1, $2, $3, $4, $5)
RETURNING
	*;

-- name: GetUser :one
SELECT
	*
FROM
	"user"
WHERE
	id = $1
LIMIT
	1;

-- name: GetUserFromEmail :one
SELECT
	*
FROM
	"user"
WHERE
	email = $1
LIMIT
	1;

-- name: GetUserFromUsername :one
SELECT
	*
FROM
	"user"
WHERE
	username = $1
LIMIT
	1;

-- name: ListUsers :many
SELECT
	*
FROM
	"user"
LIMIT
	$1
OFFSET
	$2;

-- name: DeleteUser :exec
DELETE FROM "user"
WHERE
	id = $1;