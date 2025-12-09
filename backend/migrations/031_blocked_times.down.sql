-- Rollback migration for blocked_times table

BEGIN;

DROP POLICY IF EXISTS blocked_times_tenant_isolation ON blocked_times;
DROP TABLE IF NOT EXISTS blocked_times CASCADE;

COMMIT;
