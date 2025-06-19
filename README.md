# Social media aggregator

Simple, containerized and self-contained social media aggregator.

### Run instructions

- Insert provided mastodon user token to `docker-compose.yml` in place of
  `INSERT_USER_TOKEN_HERE`

- Make sure you have [Docker](https://docs.docker.com/engine/install/) installed

- In the root of the project run:

  ```
  docker compose up
  ```

### Notes on technical decisions

- Mosquitto MQTT

  This project is using Mosquitto MQTT for decoupling specific social media
  aggregators and the common processor and client-facing API. In the real-world
  scenario I would recommend using Google Cloud Pub/Sub for its scalability and
  ease of use, but for the sake of this self-contained test assignment I opted
  for open-source Mosquitto MQTT as a substitute.

- SQLite

  SQLite is used as a simple, self-contained SQL database for this project, in
  the real-world project it would be an actual cloud hosted SQL database.

- techhub.social

  Instead of `mastodon.social` I used `techhub.social` - one of the mastodon
  websites in the "federation", because I was unable to register a new account
  with `mastodon.social`, it was not sending confirmation emails at the time,
  presumably due to high load. All of the documentation that is related to
  `mastodon.social`, also true for `techhub.social`, so I assumed them to be
  equal for this project.
