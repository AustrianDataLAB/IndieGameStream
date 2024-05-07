CREATE TABLE IF NOT EXISTS db_state (
    migrations int
);

CREATE TABLE IF NOT EXISTS games (
     ID varchar(36) NOT NULL primary key,
     Title varchar(255),
     StorageLocation varchar(255),
     Status varchar(255),
     Url varchar(255)
);

INSERT INTO db_state VALUES (0);