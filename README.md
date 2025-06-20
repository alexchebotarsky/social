# Social media aggregator

Simple, containerized and self-contained social media aggregator.

### Run instructions

- Insert provided Mastodon user token to `docker-compose.yml` in place of
  `INSERT_USER_TOKEN_HERE`

- Make sure you have [Docker](https://docs.docker.com/engine/install/) installed

- In the root of the project run:

  ```
  docker compose up
  ```

### Notes on technical decisions

- Docker Compose

  In order for this project to be easily and predictably executable I made use
  of `docker-compose.yml`. Though, this would not be a good choice for a
  real-world application, each folder would most likely be a separate deployment
  in the cloud infrastructure.

- Mosquitto MQTT

  This project is using Mosquitto MQTT for decoupling specific social media
  ingestors and the common client-facing aggregator API. In the real-world
  scenario I would recommend using Google Cloud Pub/Sub for its scalability and
  ease of use, but for the sake of this self-contained test assignment I opted
  for open-source Mosquitto MQTT as a substitute.

- SQLite

  SQLite is used as a simple, self-contained SQL database for this project, in
  the real-world project it would be a cloud hosted SQL database instance.

- techhub.social

  Instead of `mastodon.social` I used `techhub.social` - one of the mastodon
  websites in the "federation". Because I was unable to register a new account
  with `mastodon.social`, since it was not sending confirmation emails at the
  time, presumably due to high load. All of the documentation that is related to
  `mastodon.social`, also true for `techhub.social`, so I assumed them to be
  equal for this project.

### OpenAPI

Main API `social-media-aggregator` has OpenAPI documentation hosted at
`localhost:8000/docs`, explaining the API for querying the aggregated posts.
