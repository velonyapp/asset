CREATE TABLE images (
    id                CHAR(36) PRIMARY KEY,
    tags              TEXT NOT NULL,
    source_object_key VARCHAR(128) NOT NULL,
    object_key        VARCHAR(128) NOT NULL,
    state             VARCHAR(32) NOT NULL,
    width             INT NOT NULL,
    height            INT NOT NULL,
    format            VARCHAR(16) NOT NULL,
    create_time       DATETIME(6) NOT NULL,
    update_time       DATETIME(6) NOT NULL,

    INDEX idx_images_update_time       (update_time),
    INDEX idx_images_source_object_key (source_object_key),
    INDEX idx_images_object_key        (object_key),
    INDEX idx_images_format            (format)
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

CREATE TABLE outbox_jobs (
    id             CHAR(36) PRIMARY KEY,
    type           TEXT NOT NULL,
    payload        JSON NOT NULL,
    occur_time     DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6)
) ENGINE = InnoDB;
