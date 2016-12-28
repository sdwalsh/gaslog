DROP SCHEMA public CASCADE;
CREATE SCHEMA public;
CREATE EXTENSION pgcrypto;

------------------------------------------------------------
------------------------------------------------------------

CREATE TYPE user_role AS ENUM ('ADMIN', 'MOD', 'MEMBER', 'BANNED');

CREATE TABLE users (
  id             uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  uname          text NOT NULL,
  role           user_role NOT NULL DEFAULT 'MEMBER',
  digest         text NOT NULL,
  email          text NOT NULL,
  last_online_at timestamptz NOT NULL DEFAULT NOW(),
  created_at     timestamptz NOT NULL DEFAULT NOW()
);

CREATE TABLE log_book (
  id              uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id         uuid NOT NULL REFERENCES user(id)
)

CREATE TABLE log_entry (
  id                    uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  log_book_id           uuid NOT NULL REFERENCES log_book(id),
  created_at            timestampz NOT NULL DEFAULT NOW(),
  miles                 bigint NOT NULL,
  cost_per_gallon       money NOT NULL,
  cost_total            money NOT NULL,
  location              text NOT NULL
)

------------------------------------------------------------
------------------------------------------------------------

CREATE TABLE sessions (
  id            uuid PRIMARY KEY,
  user_id       uuid NOT NULL REFERENCES users(id),
  ip_address    inet NOT NULL,
  user_agent    text NULL,
  logged_out_at timestamptz NULL,
  expired_at    timestamptz NOT NULL DEFAULT NOW() + INTERVAL '2 weeks',
  created_at    timestamptz NOT NULL DEFAULT NOW()
);

-- Speed up user_id FK joins
CREATE INDEX sessions__user_id ON sessions (user_id);

CREATE VIEW active_sessions AS
  SELECT *
  FROM sessions
  WHERE expired_at > NOW()
    AND logged_out_at IS NULL
;

------------------------------------------------------------
------------------------------------------------------------

CREATE OR REPLACE FUNCTION ip_root(ip_address inet) RETURNS inet AS
$$
  DECLARE
    masklen int;
  BEGIN
    masklen := CASE family(ip_address) WHEN 4 THEN 24 ELSE 48 END;
    RETURN host(network(set_masklen(ip_address, masklen)));
  END;
$$ LANGUAGE plpgsql IMMUTABLE;

CREATE TABLE ratelimits (
  id             bigserial        PRIMARY KEY,
  ip_address     inet             NOT NULL,
  created_at     timestamptz      NOT NULL DEFAULT NOW()
);

CREATE INDEX ratelimits__ip_root ON ratelimits (ip_root(ip_address));
