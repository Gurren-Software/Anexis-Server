-- Create "users" table
CREATE TABLE "public"."users" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "email" text NOT NULL,
  "name" text NOT NULL,
  "password_hash" text NOT NULL,
  "age" smallint NULL,
  "birthday" timestamptz NULL,
  "member_number" text NULL,
  "activated_at" timestamptz NULL,
  "storage_quota" bigint NULL DEFAULT 5368709120,
  "storage_used" bigint NULL DEFAULT 0,
  PRIMARY KEY ("id")
);
-- Create index "idx_users_deleted_at" to table: "users"
CREATE INDEX "idx_users_deleted_at" ON "public"."users" ("deleted_at");
-- Create index "idx_users_email" to table: "users"
CREATE UNIQUE INDEX "idx_users_email" ON "public"."users" ("email");
-- Create index "idx_users_member_number" to table: "users"
CREATE UNIQUE INDEX "idx_users_member_number" ON "public"."users" ("member_number");
-- Create "backup_jobs" table
CREATE TABLE "public"."backup_jobs" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "user_id" uuid NOT NULL,
  "type" text NOT NULL,
  "status" text NULL DEFAULT 'pending',
  "archive_key" text NULL,
  "archive_size" bigint NULL,
  "total_files" bigint NULL DEFAULT 0,
  "processed_files" bigint NULL DEFAULT 0,
  "total_bytes" bigint NULL DEFAULT 0,
  "processed_bytes" bigint NULL DEFAULT 0,
  "started_at" timestamptz NULL,
  "completed_at" timestamptz NULL,
  "expires_at" timestamptz NULL,
  "last_error" text NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_backup_jobs_user" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create index "idx_backup_jobs_deleted_at" to table: "backup_jobs"
CREATE INDEX "idx_backup_jobs_deleted_at" ON "public"."backup_jobs" ("deleted_at");
-- Create index "idx_backup_jobs_user_id" to table: "backup_jobs"
CREATE INDEX "idx_backup_jobs_user_id" ON "public"."backup_jobs" ("user_id");
-- Create "files" table
CREATE TABLE "public"."files" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "user_id" uuid NOT NULL,
  "name" text NOT NULL,
  "original_name" text NOT NULL,
  "mime_type" text NOT NULL,
  "size" bigint NOT NULL,
  "storage_key" text NOT NULL,
  "storage_path" text NOT NULL,
  "checksum" text NOT NULL,
  "status" text NULL DEFAULT 'pending',
  "is_encrypted" boolean NULL DEFAULT false,
  "is_compressed" boolean NULL DEFAULT false,
  "parent_id" uuid NULL,
  "description" text NULL,
  "tags" text NULL,
  "uploaded_at" timestamptz NULL,
  "processed_at" timestamptz NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_files_parent" FOREIGN KEY ("parent_id") REFERENCES "public"."files" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "fk_users_files" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create index "idx_files_deleted_at" to table: "files"
CREATE INDEX "idx_files_deleted_at" ON "public"."files" ("deleted_at");
-- Create index "idx_files_parent_id" to table: "files"
CREATE INDEX "idx_files_parent_id" ON "public"."files" ("parent_id");
-- Create index "idx_files_storage_key" to table: "files"
CREATE UNIQUE INDEX "idx_files_storage_key" ON "public"."files" ("storage_key");
-- Create index "idx_files_user_id" to table: "files"
CREATE INDEX "idx_files_user_id" ON "public"."files" ("user_id");
-- Create "links" table
CREATE TABLE "public"."links" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "user_id" uuid NOT NULL,
  "file_id" uuid NOT NULL,
  "token" text NOT NULL,
  "type" text NOT NULL,
  "access_type" text NULL DEFAULT 'public',
  "password" text NULL,
  "max_downloads" bigint NULL,
  "download_count" bigint NULL DEFAULT 0,
  "allowed_emails" text NULL,
  "expires_at" timestamptz NULL,
  "last_accessed_at" timestamptz NULL,
  "name" text NULL,
  "description" text NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_files_links" FOREIGN KEY ("file_id") REFERENCES "public"."files" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "fk_users_links" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create index "idx_links_deleted_at" to table: "links"
CREATE INDEX "idx_links_deleted_at" ON "public"."links" ("deleted_at");
-- Create index "idx_links_file_id" to table: "links"
CREATE INDEX "idx_links_file_id" ON "public"."links" ("file_id");
-- Create index "idx_links_token" to table: "links"
CREATE UNIQUE INDEX "idx_links_token" ON "public"."links" ("token");
-- Create index "idx_links_user_id" to table: "links"
CREATE INDEX "idx_links_user_id" ON "public"."links" ("user_id");
-- Create "migration_jobs" table
CREATE TABLE "public"."migration_jobs" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "user_id" uuid NOT NULL,
  "provider" text NOT NULL,
  "status" text NULL DEFAULT 'pending',
  "total_files" bigint NULL DEFAULT 0,
  "processed_files" bigint NULL DEFAULT 0,
  "failed_files" bigint NULL DEFAULT 0,
  "total_bytes" bigint NULL DEFAULT 0,
  "processed_bytes" bigint NULL DEFAULT 0,
  "started_at" timestamptz NULL,
  "completed_at" timestamptz NULL,
  "last_error" text NULL,
  "access_token" text NULL,
  "refresh_token" text NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_migration_jobs_user" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create index "idx_migration_jobs_deleted_at" to table: "migration_jobs"
CREATE INDEX "idx_migration_jobs_deleted_at" ON "public"."migration_jobs" ("deleted_at");
-- Create index "idx_migration_jobs_user_id" to table: "migration_jobs"
CREATE INDEX "idx_migration_jobs_user_id" ON "public"."migration_jobs" ("user_id");
