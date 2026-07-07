CREATE TABLE groups (
  id VARCHAR(32) PRIMARY KEY,
  name VARCHAR(32) NOT NULL,
  created_at TIMESTAMP DEFAULT now(),
  updated_at TIMESTAMP DEFAULT now(),
  deleted_at TIMESTAMP DEFAULT NULL
);
