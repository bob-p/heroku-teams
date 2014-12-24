Simple web interface for managing teams on heroku

psql -f database.sql

cp .env.example to .env and add relevent info

docker build -t heroku-teams .

docker run -d -p 8080:3000 --name myapp heroku-teams

goconvey
