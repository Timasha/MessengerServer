CREATE TABLE IF NOT EXISTS users(
    login text PRIMARY KEY,
    password text,
    refreshBodies text[]
);