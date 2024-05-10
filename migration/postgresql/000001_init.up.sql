CREATE TABLE IF NOT EXISTS instrument
(
    id         INT,
    changed    INT,
    title      TEXT,
    short_title TEXT,
    symbol     TEXT,
    rate       FLOAT,
    PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS company
(
    id        INT,
    changed   INT,
    title     TEXT,
    full_title TEXT,
    www       TEXT,
    country   INT,
    PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS "user"
(
    id       INT,
    changed  INT,
    login    TEXT,
    currency INT,
    parent   INT,
    PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS country
(
    id       INT,
    title    TEXT,
    currency INT,
    domain   TEXT,
    PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS account
(
    id                    UUID,
    changed               INT,
    "user"                INT,
    role                  INT,
    instrument            INT,
    company               INT,
    type                  TEXT,
    title                 TEXT,
    sync_id               TEXT[],
    balance               FLOAT,
    start_balance         FLOAT,
    credit_limit          FLOAT,
    in_balance            INT,
    savings               BOOLEAN,
    enable_correction     INT,
    enable_sms            INT,
    archive               INT,
    capitalization        BOOLEAN,
    percent               FLOAT,
    start_date            TEXT,
    end_date_offset       INT,
    end_date_offset_interval TEXT,
    payoff_step           INT,
    payoff_interval       TEXT,
    PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS tag
(
    id            UUID,
    changed       INT,
    "user"        INT,
    title         TEXT,
    parent        TEXT,
    icon          TEXT,
    picture       TEXT,
    color         BIGINT,
    show_income   INT,
    show_outcome  INT,
    budget_income INT,
    budget_outcome INT,
    required      BOOLEAN,
    PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS merchant
(
    id      UUID,
    changed INT,
    "user"  INT,
    title   TEXT,
    PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS reminder
(
    id                UUID,
    changed           INT,
    "user"            INT,
    income_instrument INT,
    income_account    TEXT,
    income            FLOAT,
    outcome_instrument INT,
    outcome_account   TEXT,
    outcome           FLOAT,
    tag               UUID[],
    merchant          UUID,
    payee             TEXT,
    comment           TEXT,
    interval          TEXT,
    step              INT,
    points            INT[],
    start_date        TEXT,
    end_date          TEXT,
    notify            INT,
    PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS reminder_marker
(
    id                UUID,
    changed           INT,
    "user"            INT,
    income_instrument INT,
    income_account    TEXT,
    income            FLOAT,
    outcome_instrument INT,
    outcome_account   TEXT,
    outcome           FLOAT,
    tag               UUID[],
    merchant          UUID,
    payee             TEXT,
    comment           TEXT,
    date              TEXT,
    reminder          UUID,
    state             TEXT,
    notify            INT,
    PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS transaction
(
    id                  UUID,
    changed             INT,
    created             INT,
    "user"              INT,
    deleted             BOOLEAN,
    hold                BOOLEAN,
    income_instrument   INT,
    income_account      TEXT,
    income              FLOAT,
    outcome_instrument  INT,
    outcome_account     TEXT,
    outcome             FLOAT,
    tag                 UUID[],
    merchant            UUID,
    payee               TEXT,
    original_payee      TEXT,
    comment             TEXT,
    date                TEXT,
    mcc                 INT,
    reminder_marker     UUID,
    op_income           FLOAT,
    op_income_instrument INT,
    op_outcome          FLOAT,
    op_outcome_instrument INT,
    latitude            FLOAT,
    longitude           FLOAT,
    PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS budget
(
    changed     INT,
    "user"      INT,
    tag         UUID,
    date        TEXT,
    income      FLOAT,
    income_lock INT,
    outcome     FLOAT,
    outcome_lock INT
);