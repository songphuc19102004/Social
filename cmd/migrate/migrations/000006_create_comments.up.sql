CREATE TABLE IF NOT EXISTS comments(
  id BIGINT GENERATED ALWAYS AS IDENTITY,
  post_id BIGINT NOT NULL,
  user_id BIGINT NOT NULL,
  content text NOT NULL,
  created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
  updated_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
  PRIMARY KEY (id),
  FOREIGN KEY (post_id) REFERENCES posts(id),
  FOREIGN KEY (user_id) REFERENCES users(id)
);
