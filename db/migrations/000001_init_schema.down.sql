-- Drop foreign keys
ALTER TABLE "logs" DROP CONSTRAINT IF EXISTS "logs_user_id_fkey";
ALTER TABLE "users" DROP CONSTRAINT IF EXISTS "users_tenant_id_fkey";

-- Drop indexes
DROP INDEX IF EXISTS "user_id_tenant_id_idx";

-- Drop tables
DROP TABLE IF EXISTS "logs";
DROP TABLE IF EXISTS "users";
DROP TABLE IF EXISTS "tenant";