ALTER TABLE users DROP CONSTRAINT IF EXISTS users_preferred_language_check;
ALTER TABLE users DROP COLUMN IF EXISTS preferred_language;
