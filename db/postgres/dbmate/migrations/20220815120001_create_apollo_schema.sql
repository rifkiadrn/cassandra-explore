-- migrate:up
SET LOCAL lock_timeout = '60s';
CREATE SCHEMA IF NOT EXISTS cassandra_users;

-- migrate:down
SET LOCAL lock_timeout = '60s';
DROP SCHEMA IF EXISTS cassandra_users;