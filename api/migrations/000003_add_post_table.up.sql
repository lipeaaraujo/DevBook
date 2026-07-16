CREATE TABLE POSTS (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  title VARCHAR(255) NOT NULL,
  description VARCHAR(255) NOT NULL,
  author_id UUID NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT now(),
  updated_at TIMESTAMP,
  FOREIGN KEY(author_id)
    REFERENCES USERS(id)
    ON DELETE CASCADE
);
