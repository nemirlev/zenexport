CREATE TABLE IF NOT EXISTS reminder_marker
(
    id                UUID,
    changed           Int32,
    user              Int32,
    income_instrument  Int32,
    income_account     String,
    income            Float64,
    outcome_instrument Int32,
    outcome_account    String,
    outcome           Float64,
    tag               Array(UUID),
    merchant          Nullable(UUID),
    payee             String,
    comment           String,
    date              String,
    reminder          UUID,
    state             String,
    notify            UInt8
) ENGINE = MergeTree PRIMARY KEY id;