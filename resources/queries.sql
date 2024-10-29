-- name: get-user-by-external-id
SELECT *
FROM auth.users
WHERE
    external_identifier = $1;

-- name: create-user
INSERT INTO auth.users(external_identifier, name, username, email)
VALUES
    ($1, $2, $3, $4);

-- name: get-services
SELECT *
FROM auth.services;

-- name: get-user-permissions
SELECT s.name, level
FROM auth.permission_assignments
JOIN auth.services s ON s.id = permission_assignments.service
WHERE user_id = $1;