//
// ATTENTION!!! DO NOT EDIT THIS FILE. IT IS CREATED FOLLOWING SCHEMA MIGRATION STRATEGY
//

CREATE TABLE IF NOT EXISTS tasks (
    id                   uuid,
    partner_id           text,
    created_at           timestamp,
    created_by           text,
    description          text,
    end_run_time         timestamp,
    external_task        boolean,
    last_run_at          timestamp,
    managed_endpoint_id  uuid,
    name                 text,
    origin_id            uuid,
    parameters           text,
    regularity           int,
    run_time             timestamp,
    run_time_unix        timestamp,
    schedule             text,
    start_run_time       timestamp,
    started_at           timestamp,
    state                int,
    target_type          int,
    trigger              text,
    type                 text,
    location             text,
    PRIMARY KEY          (partner_id, external_task, managed_endpoint_id, id)
);

CREATE MATERIALIZED VIEW IF NOT EXISTS tasks_by_runtime AS
    SELECT *
    FROM tasks
    WHERE run_time_unix IS NOT NULL AND partner_id IS NOT NULL AND id IS NOT NULL AND managed_endpoint_id IS NOT NULL AND external_task IS NOT NULL
    PRIMARY KEY (run_time_unix, id, partner_id, managed_endpoint_id, external_task);
