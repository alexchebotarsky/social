# Social media aggregator

Simple, containerized and self-contained social media aggregator.

### Run instructions

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
