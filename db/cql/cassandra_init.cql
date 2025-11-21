CREATE KEYSPACE IF NOT EXISTS blogs
WITH replication = {'class': 'SimpleStrategy', 'replication_factor': 2};

CREATE TABLE IF NOT EXISTS blogs.users (
  id uuid PRIMARY KEY,
  name text,
  username text
);

CREATE TABLE blogs.blogs_by_author (
    author_id uuid,
    username text,
    ts timeuuid,
    id uuid,
    content text,
    PRIMARY KEY (author_id, ts)
) WITH CLUSTERING ORDER BY (ts DESC);