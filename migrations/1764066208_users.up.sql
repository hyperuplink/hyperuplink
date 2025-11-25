CREATE TABLE users (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  username VARCHAR(32),
  role VARCHAR(16),
  password VARCHAR(32),
  password_reset_token VARCHAR(32),
  password_reset_sent_at TIMESTAMP,
  email VARCHAR(254),
  email_unconfirmed VARCHAR(254),
  email_confirmation_token VARCHAR(32),
  email_confirmation_sent_at TIMESTAMP,
  email_confirmed_at TIMESTAMP,
  otp_enabled BOOLEAN,
  otp_secret VARCHAR(32),
  otp_timestep SMALLINT,
  sign_in_last_at TIMESTAMP,
  sign_in_failed_attempts SMALLINT,
  sign_in_locked_at TIMESTAMP,
  sign_in_unlock_token VARCHAR(32),
  profile_picture VARCHAR(32),
  signature VARCHAR(256),
  created_at TIMESTAMP,
  updated_at TIMESTAMP,
  banned_at TIMESTAMP,
  deleted_at TIMESTAMP
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_users_email ON users (
  email, deleted_at
);
CREATE UNIQUE INDEX IF NOT EXISTS idx_users_username ON users (
  username, deleted_at
);
