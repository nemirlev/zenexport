CREATE TABLE IF NOT EXISTS instrument
(
    id         Int32,
    changed    Int32,
    title      String,
    short_title String,
    symbol     String,
    rate       Float64
) ENGINE = MergeTree PRIMARY KEY id;