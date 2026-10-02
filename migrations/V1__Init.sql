CREATE TABLE images (
    id            CHAR(36) PRIMARY KEY,
    tags          JSON NOT NULL,
    object_key    VARCHAR(128) NOT NULL UNIQUE,
    object_exists BOOLEAN NOT NULL,
    create_time   TIMESTAMP(6) NOT NULL
) ENGINE = InnoDB;

CREATE TABLE outbox_events (
    id             CHAR(36) PRIMARY KEY,
    aggregate_id   TEXT NOT NULL,
    aggregate_type TEXT NOT NULL,
    type           TEXT NOT NULL,
    payload        JSON NOT NULL,
    occur_time     TIMESTAMP(6) NOT NULL
) ENGINE = InnoDB;
