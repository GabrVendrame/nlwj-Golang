-- Write your migrate up statements here

CREATE TABLE IF NOT EXISTS links (
    "id"        uuid    PRIMARY KEY     NOT NULL    DEFAULT gen_random_uuid(),
    "trip_id"   uuid                    NOT NULL,
    "title"     VARCHAR(255)            NOT NULL,
    "url"       VARCHAR(255)            NOT NULL,

    FOREIGN KEY (trip_id) REFERENCES trips(id)
        ON UPDATE CASCADE
        ON DELETE CASCADE
);

---- create above / drop below ----

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.

DROP TABLE IF EXISTS links;
