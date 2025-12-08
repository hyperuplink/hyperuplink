CREATE TABLE users (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  username VARCHAR(32),
  role VARCHAR(16) DEFAULT 'user',
  member_of VARCHAR(32) [] DEFAULT NULL,
  password VARCHAR(255),
  password_reset_token VARCHAR(32) DEFAULT '',
  password_reset_sent_at TIMESTAMP DEFAULT NULL,
  email VARCHAR(254),
  email_unconfirmed VARCHAR(254),
  email_confirmation_token VARCHAR(32) DEFAULT '',
  email_confirmation_sent_at TIMESTAMP DEFAULT NULL,
  email_confirmed_at TIMESTAMP DEFAULT NULL,
  language VARCHAR(2) DEFAULT 'en',
  otp_enabled BOOLEAN DEFAULT FALSE,
  otp_secret VARCHAR(32) DEFAULT '',
  otp_timestep SMALLINT DEFAULT 0,
  sign_in_last_at TIMESTAMP DEFAULT NULL,
  sign_in_failed_attempts SMALLINT DEFAULT 0,
  sign_in_locked_at TIMESTAMP DEFAULT NULL,
  sign_in_unlock_token VARCHAR(32) DEFAULT '',
  profile_picture VARCHAR(32) DEFAULT '',
  signature VARCHAR(256) DEFAULT '',
  created_at TIMESTAMP DEFAULT now(),
  updated_at TIMESTAMP DEFAULT now(),
  confirmed_at TIMESTAMP DEFAULT NULL,
  banned_at TIMESTAMP DEFAULT NULL,
  deleted_at TIMESTAMP DEFAULT NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_users_email ON users (
  email, deleted_at
);
CREATE UNIQUE INDEX IF NOT EXISTS idx_users_username ON users (
  username, deleted_at
);
