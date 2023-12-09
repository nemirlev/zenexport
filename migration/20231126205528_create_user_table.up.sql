CREATE TABLE IF NOT EXISTS user
(
    id       Int32,
    changed  Int32,
    login    Nullable(String),
    currency Int32,
    parent   Nullable(Int32)
) ENGINE = MergeTree PRIMARY KEY id;