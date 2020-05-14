ALTER TABLE tasks DROP credentials;
ALTER TABLE task_by_runtime_unix_mv DROP credentials;
ALTER TABLE tasks_order_by_last_task_instance_id_mv DROP credentials;
ALTER TABLE tasks_by_runtime_mv DROP credentials;
ALTER TABLE tasks_by_id_mv DROP credentials;
DROP TYPE credentials;