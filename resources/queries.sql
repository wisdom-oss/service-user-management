-- USER-RELATED QUERIES --

-- name: get-user-by-external-id
SELECT *
FROM auth.users
WHERE
    external_identifier = $1;

-- name: get-user-by-internal-id
SELECT *
FROM auth.users
WHERE
    id = $1::uuid;

-- name: get-users
SELECT *
FROM auth.users;

-- name: create-user
INSERT INTO auth.users(external_identifier, name, username, email)
VALUES
    ($1, $2, $3, $4);

-- name: delete-user
DELETE FROM auth.users
WHERE id = $1::uuid;

-- TOKEN RELATED QUERIES --

-- name: check-for-refresh-token
SELECT EXISTS(SELECT id
              FROM auth.refresh_tokens
              WHERE
                  id = $1
                  AND active IS TRUE
                  AND expires_at > NOW());

-- name: register-refresh-token
INSERT INTO auth.refresh_tokens(id, active, expires_at)
VALUES
    ($1, TRUE, $2);


-- name: revoke-refresh-token
UPDATE auth.refresh_tokens
SET active = FALSE
WHERE
    id = $1;

-- name: cleanup-expired-tokens
DELETE
FROM auth.refresh_tokens
WHERE
    expires_at < NOW() OR active is not TRUE;


-- SERVICE RELATED QUERIES --

-- name: get-services
SELECT *
FROM auth.services;

-- name: get-service-by-internal-id
SELECT *
FROM auth.services
WHERE
    id = $1::uuid
LIMIT 1;

-- name: get-service-by-external-id
SELECT *
FROM auth.services
WHERE
    name = $1
LIMIT 1;

-- PERMISSION RELATED QUERIES --

-- name: get-user-permissions
SELECT s.name, level
FROM auth.permission_assignments
         JOIN auth.services s ON s.id = permission_assignments.service
WHERE
    user_id = $1;

-- name: assign-permission
INSERT INTO auth.permission_assignments(user_id, service, level)
VALUES
    ($1, $2, $3)