//
// ATTENTION!!! DO NOT EDIT THIS FILE. IT IS CREATED FOLLOWING SCHEMA MIGRATION STRATEGY
//

ALTER TABLE tasks ADD modified_at timestamp;
ALTER TABLE tasks ADD modified_by text;