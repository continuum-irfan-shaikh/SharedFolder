//
// ATTENTION!!! DO NOT EDIT THIS FILE. IT IS CREATED FOLLOWING SCHEMA MIGRATION STRATEGY
//

DROP MATERIALIZED VIEW IF EXISTS task_definitions_by_deleted;

ALTER TABLE task_definitions DROP deleted;