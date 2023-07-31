CREATE SCHEMA IF NOT EXISTS __pubsub;
CREATE TABLE IF NOT EXISTS __pubsub.stats(
    id varchar(64) NOT NULL PRIMARY KEY,
    timestamp TIMESTAMP
)
