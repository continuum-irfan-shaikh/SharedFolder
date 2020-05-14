CREATE TABLE IF NOT EXISTS platform_tasking_db.active_triggers (
    partner_id text,
    type       text,
    task_id    uuid,
    PRIMARY KEY (partner_id, type, task_id)
);