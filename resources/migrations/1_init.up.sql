BEGIN;
CREATE SCHEMA "auth";

CREATE TABLE IF NOT EXISTS auth.users
(
    -- id is used to keep track of the user in the service and it's tables.
    -- as it creates a new unique key for a user instead of directly using the
    -- subject identifier
    id                  uuid PRIMARY KEY DEFAULT gen_random_uuid(),

    -- external_identifier contains the 'sub' claim by which the user can be
    -- identified if an id token is exchanged for an access token
    external_identifier text NOT NULL,

    -- name is automatically derived from the id token and the userinfo endpoint
    -- if possible and updated if it changes
    name                text
);

CREATE TYPE auth.scope_level AS ENUM ('read', 'write', 'delete', 'admin');

CREATE TABLE auth.services
(
    id                     uuid PRIMARY KEY            DEFAULT gen_random_uuid(),
    name                   text               NOT NULL,
    description            text,
    supported_scope_levels auth.scope_level[] NOT NULL DEFAULT ('read', 'write', 'delete', 'admin')
);

CREATE TABLE IF NOT EXISTS auth.permission_assignments
(
    assignment_id bigserial,
    user_id uuid REFERENCES auth.users (id) MATCH FULL ON DELETE CASCADE ON UPDATE RESTRICT,
    service uuid REFERENCES auth.services (id) MATCH FULL ON DELETE CASCADE ON UPDATE CASCADE,
    level   auth.scope_level NOT NULL
);

CREATE TABLE IF NOT EXISTS auth.access_tokens
(
    id         uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    subject    uuid REFERENCES auth.users (id),
    issued     timestamp        DEFAULT NOW()::timestamp,
    expires_at timestamp GENERATED ALWAYS AS ( issued + INTERVAL '1 hour' ) STORED,
    scopes     bigint[]
);

CREATE TABLE IF NOT EXISTS auth.refresh_token
(
    id         uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    subject    uuid REFERENCES auth.users (id),
    issued     timestamp        DEFAULT NOW()::timestamp,
    expires_at timestamp GENERATED ALWAYS AS ( issued + INTERVAL '1 day' ) STORED,
    scopes     bigint[]
);

CREATE OR REPLACE FUNCTION drop_old_access_token() RETURNS void AS
$$
BEGIN
    DELETE FROM auth.access_tokens WHERE expires_at < NOW();
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION drop_old_refresh_token() RETURNS void AS
$$
BEGIN
    DELETE FROM auth.refresh_token WHERE expires_at < NOW();
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER remove_old_access_tokens
    BEFORE INSERT
    ON auth.access_tokens
    FOR EACH STATEMENT
EXECUTE FUNCTION drop_old_access_token();

CREATE TRIGGER remove_old_refresh_tokens
    BEFORE INSERT
    ON auth.refresh_token
    FOR EACH STATEMENT
EXECUTE FUNCTION drop_old_refresh_token();

