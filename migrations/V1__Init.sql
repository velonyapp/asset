CREATE TABLE images (
    id            CHAR(36) PRIMARY KEY,
    tags          TEXT NOT NULL,
    object_key    VARCHAR(128) NOT NULL,
    object_exists BOOLEAN NOT NULL,
    create_time   TIMESTAMP(6) NOT NULL,

    UNIQUE KEY uq_images_object_key (object_key)
) ENGINE = InnoDB;

CREATE TABLE outbox_events (
    id             CHAR(36) PRIMARY KEY,
    type           TEXT NOT NULL,
    aggregate_id   TEXT NOT NULL,
    aggregate_type TEXT NOT NULL,
    tags           TEXT NOT NULL,
    payload        JSON NOT NULL,
    occur_time     TIMESTAMP(6) NOT NULL
) ENGINE = InnoDB;
