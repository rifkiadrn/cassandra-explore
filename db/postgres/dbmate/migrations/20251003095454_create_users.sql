-- migrate:up
CREATE TABLE cassandra_users.users (
    id UUID NOT NULL PRIMARY KEY DEFAULT uuid_generate_v4(),
    created_at BIGINT NOT NULL DEFAULT DATE_PART('EPOCH', NOW()),
    updated_at BIGINT NOT NULL DEFAULT DATE_PART('EPOCH', NOW()),
    name VARCHAR(255) NOT NULL,
    username VARCHAR(255) NOT NULL,
    password VARCHAR(255) NOT NULL,
    token VARCHAR(255) NOT NULL
);

-- migrate:down
DROP TABLE IF EXISTS cassandra_users.users;

