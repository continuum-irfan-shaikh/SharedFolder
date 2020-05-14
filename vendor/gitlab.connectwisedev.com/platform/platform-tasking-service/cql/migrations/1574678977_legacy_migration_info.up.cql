CREATE TABLE IF NOT EXISTS script_migration_info
(
    partner_id         text,
    legacy_id          text,
    legacy_template_id text,
    is_sequence        boolean,
    is_parametrized    boolean,
    origin_id          text,
    definition_id      text,
    reason             text,
    PRIMARY KEY (partner_id, legacy_id)
);

CREATE TABLE IF NOT EXISTS job_migration_info
(
    partner_id         text,
    legacy_job_id      text,
    legacy_script_id   text,
    legacy_template_id text,
    origin_id          text,
    definition_id      text,
    type               text,
    task_id            text,
    reason             text,
    PRIMARY KEY (partner_id, legacy_job_id)
);