CREATE TABLE IF NOT EXISTS tag
(
    id            UUID,
    changed       Int32,
    user          Int32,
    title         String,
    parent        Nullable(String),
    icon          Nullable(String),
    picture       Nullable(String),
    color         Nullable(Int32),
    show_income    UInt8,
    show_outcome   UInt8,
    budget_income  UInt8,
    budget_outcome UInt8,
    required      Nullable(BOOL)
) ENGINE = MergeTree PRIMARY KEY id;