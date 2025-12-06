DROP TABLE IF EXISTS learning_materials;
DROP TABLE IF EXISTS enrollments;
DROP TABLE IF EXISTS courses;
DROP TABLE IF EXISTS users;

DROP TYPE IF EXISTS material_type;
DROP TYPE IF EXISTS user_status;
DROP TYPE IF EXISTS user_role;

-- remove the avatar_url column from the users table
ALTER TABLE users
DROP COLUMN avatar_url;