# Osusu

## About

Osusu is an app for getting recommendations on what meals to eat in a group based on the ratings of each member of the group, and the cost, effort, healthiness, and recency of the meal. Osusu is an installable [Progressive Web Application (PWA)](https://web.dev/progressive-web-apps/) written in Go and made using [Go-App](https://github.com/maxence-charriere/go-app). [Osusu is hosted here](https://osusu.fly.dev).

## Run Locally

To run Osusu locally, first get the source code using `git clone https://github.com/kplat1/osusu`. Then, navigate to the directory of the app and build and run it by running `make`. The app will now be available at [localhost:8000](http://localhost:8000). Thirdly, you need to locally host a PostgreSQL database on your device. To do this, [download PostgreSQL](https://www.postgresql.org/download/) if you don't already have it installed. Then, set up a database by following the steps in the installer. Finally, set the DATABASE_URL environment variable to the URL of this database using `export DATABASE_URL=postgres://postgres:{database_server_password}@localhost:{database_server_port}/{database_name}?sslmode=disable`. You will have set the database server password when setting it up, the database server port is likely 5432 unless you set it to something else, and the database name is likely postgres unless you set it to something else.

## Deploy

If you have access to the fly.io account hosting the app, you can deploy the app by running `fly deploy`.