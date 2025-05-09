BEGIN TRANSACTION;

CREATE EXTENSION pgcrypto;

CREATE TABLE users
(
    id       SERIAL PRIMARY KEY,
    login    VARCHAR(32)  NOT NULL UNIQUE,
    password VARCHAR(128) NOT NULL
);

CREATE TABLE data
(
    id       SERIAL PRIMARY KEY,
    owner    SERIAL       NOT NULL,
    data     BYTEA        NOT NULL,
    metadata VARCHAR(128) NULL
);

COMMIT;