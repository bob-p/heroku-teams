CREATE database heroku_teams_dev;

\c heroku_teams_dev;

CREATE table IF NOT EXISTS users(
  id SERIAL PRIMARY KEY,
  email varchar(256),
  created_at timestamp
);
