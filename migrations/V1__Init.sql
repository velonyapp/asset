CREATE TABLE images (
    id                CHAR(36) PRIMARY KEY,
    tags              TEXT NOT NULL,
    source_object_key VARCHAR(128) NOT NULL,
    object_key        VARCHAR(128) NOT NULL,
    state             VARCHAR(32) NOT NULL,
    create_time       DATETIME(6) NOT NULL,
    delete_time       DATETIME(6),

    INDEX idx_images_delete_time       (delete_time),
    INDEX idx_images_source_object_key (source_object_key),
    INDEX idx_images_object_key        (object_key)
) ENGINE = InnoDB;

CREATE TABLE outbox_events (
    id             CHAR(36) PRIMARY KEY,
    type           TEXT NOT NULL,
    aggregate_id   TEXT NOT NULL,
    aggregate_type TEXT NOT NULL,
    tags           TEXT NOT NULL,
    payload        JSON NOT NULL,
    occur_time     DATETIME(6) NOT NULL
) ENGINE = InnoDB;
