-- Cleanup script: normalize empty strings in unique columns to NULL to avoid unique constraint conflicts
-- Users: phone/email unique columns
UPDATE users SET phone = NULL WHERE phone = '';
UPDATE users SET email = NULL WHERE email = '';

