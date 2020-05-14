CREATE TABLE IF NOT EXISTS trigger_definitions
(
    id          TEXT,
    data        TEXT,
    name        TEXT,
    description TEXT,
    PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS trigger_policy_counter
(
    trigger_id TEXT,
    policy_id  TEXT,
    count      COUNTER,
    PRIMARY KEY (trigger_id, policy_id)
);