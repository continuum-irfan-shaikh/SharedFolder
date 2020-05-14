CREATE MATERIALIZED VIEW IF NOT EXISTS task_definitions_by_deleted AS
    SELECT *
    FROM task_definitions
    WHERE partner_id IS NOT NULL AND id IS NOT NULL AND deleted IS NOT NULL
    PRIMARY KEY (partner_id, deleted, id);