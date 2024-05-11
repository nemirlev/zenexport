CREATE TABLE IF NOT EXISTS company
(
    id        Int32,
    changed   Int32,
    title     String,
    full_title String,
    www       String,
    country   Int32
) ENGINE = MergeTree PRIMARY KEY id;