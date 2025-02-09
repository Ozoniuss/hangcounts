-- Drop indexes first to avoid dependency issues
DROP INDEX IF EXISTS idx_hangout_individuals_created_at;
DROP INDEX IF EXISTS idx_hangout_individuals_individual;
DROP INDEX IF EXISTS idx_hangout_individuals_hangout;
DROP INDEX IF EXISTS idx_individuals_email;
DROP INDEX IF EXISTS idx_individuals_name;

-- Drop tables in reverse order to respect foreign key constraints
DROP TABLE IF EXISTS hangout_individuals;
DROP TABLE IF EXISTS hangouts;
DROP TABLE IF EXISTS individuals;
