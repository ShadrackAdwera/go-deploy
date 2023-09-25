-- Migrate Up SQL
ALTER TABLE "logs"
ALTER COLUMN "description" TYPE varchar(200);
