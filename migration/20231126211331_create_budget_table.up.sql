CREATE TABLE IF NOT EXISTS budget
(
    changed     Int32,
    user        Int32,
    tag         Nullable(UUID),
    date        String,
    income      Float64,
    income_lock  UInt8,
    outcome     Float64,
    outcome_lock UInt8
) ENGINE = MergeTree() ORDER BY date;