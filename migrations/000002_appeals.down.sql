DROP TABLE IF EXISTS appeal_status_history;
DROP TABLE IF EXISTS appeal_files;
DROP TABLE IF EXISTS appeals;
DROP TABLE IF EXISTS appeal_statuses;
DROP TABLE IF EXISTS dictionaries;

ALTER TABLE users DROP COLUMN IF EXISTS role;
