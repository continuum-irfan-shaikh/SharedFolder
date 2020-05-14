
ALTER TABLE platform_tasking_db.task_instances ADD device_statuses frozen <map<UUID,int>>;
ALTER TABLE platform_tasking_db.task_instances ADD failure_count int;
ALTER TABLE platform_tasking_db.task_instances ADD success_count int;
ALTER TABLE platform_tasking_db.task_instances ADD overall_status int;