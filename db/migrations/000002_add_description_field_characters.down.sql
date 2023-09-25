-- Migrate Down SQL
ALTER TABLE "logs"
ALTER COLUMN "description" TYPE varchar(40);
