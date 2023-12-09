CREATE TABLE IF NOT EXISTS country
(
    id       Int32,
    title    String,
    currency Int32,
    domain   String
) ENGINE = MergeTree PRIMARY KEY id;