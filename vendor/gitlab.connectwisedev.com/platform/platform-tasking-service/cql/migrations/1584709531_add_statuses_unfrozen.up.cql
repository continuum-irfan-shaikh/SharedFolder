//
// ATTENTION!!! DO NOT EDIT THIS FILE. IT IS CREATED FOLLOWING SCHEMA MIGRATION STRATEGY
//

ALTER TABLE task_instances
    ADD endpoints map<uuid, int>;
ALTER TABLE task_instances_by_id_mv
    ADD endpoints map<uuid, int>;
ALTER TABLE task_instances_started_at_mv
    ADD endpoints map<uuid, int>;
