CREATE TABLE IF NOT EXISTS merchant
(
    id      UUID,
    changed Int32,
    user    Int32,
    title   String
) ENGINE = MergeTree PRIMARY KEY id;