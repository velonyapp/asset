CREATE TABLE images (
    id          CHAR(36) PRIMARY KEY,
    storage_key VARCHAR(128) NOT NULL UNIQUE,
    ready       BOOLEAN NOT NULL,
    create_time TIMESTAMP(6) NOT NULL,
    delete_time TIMESTAMP(6)
) ENGINE = InnoDB;

CREATE TABLE outbox_events (
    id             CHAR(36) PRIMARY KEY,
    aggregate_id   TEXT NOT NULL,
    aggregate_type TEXT NOT NULL,
    type           TEXT NOT NULL,
    payload        JSON NOT NULL,
    occur_time     TIMESTAMP(6) NOT NULL
) ENGINE = InnoDB;
