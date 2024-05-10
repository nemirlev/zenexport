CREATE TABLE IF NOT EXISTS reminder
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
    interval          Nullable(String),
    step              Nullable(Int32),
    points            Array(Int32),
    start_date         String,
    end_date           Nullable(String),
    notify            UInt8
) ENGINE = MergeTree PRIMARY KEY id;