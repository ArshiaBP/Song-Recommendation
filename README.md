# Song Recommendation

You send a song file and your email, and a few moments later you get back 5 similar songs picked by Spotify. Built as three small Go microservices talking through RabbitMQ, instead of one monolith doing the whole pipeline synchronously.

## How it works

1. **Request Registration Service** — `POST /send-song/:email` with an mp3 file. It stores the file in object storage (S3-compatible), creates a row in Postgres with status `pending`, and publishes the request ID to RabbitMQ.
2. **Song ID Identification Service** — consumes the queue, downloads the file, sends it to the Shazam API to get the song title, then searches Spotify for that title to get its Spotify track ID. Saves the ID and flips the status to `ready` (or `failure` if any step breaks).
3. **Email Service** — polls the database every 10 seconds for `ready`/`failure` requests. For `ready`, it asks Spotify's recommendation API for 5 similar tracks and emails them to the user via Mailgun; for `failure`, it sends an apology email instead. Either way, the request row is cleaned up afterward.

Each service only does one job and only talks to the next one through the queue or the shared database — nothing calls another service's API directly.

## Services

| Service | Responsibility | Talks to |
|---|---|---|
| `request-registeration-service` | Accepts uploads, stores file, queues request | Postgres, object storage, RabbitMQ |
| `songID-identification-service` | Identifies the song and its Spotify ID | RabbitMQ, object storage, Shazam API, Spotify API |
| `email-service` | Gets recommendations and emails the user | Postgres, Spotify API, Mailgun |

## Stack

Go, [Echo](https://echo.labstack.com/), GORM + PostgreSQL, RabbitMQ, S3-compatible object storage, Shazam & Spotify (via RapidAPI), Mailgun. Everything runs through `docker-compose`.

## Running project

```bash
docker-compose up --build
```

Each service reads its config from its own env file, referenced in `docker-compose.yml`:
- `request-registeration-service.env`, `songID-identification-service.env`, `email-service.env`

Between them you'll need to provide:
- Postgres: `DATABASE_HOST`, `DATABASE_PORT`, `DATABASE_USERNAME`, `DATABASE_PASSWORD`, `DATABASE_DB`
- RabbitMQ: `RabbitMQ_URL`, `RabbitMQ_User`, `RabbitMQ_Pass`
- Object storage: `ACCESS_KEY`, `SECRET_KEY`, `ENDPOINT`, `BUCKET_NAME`
- RapidAPI (Shazam + Spotify): `API_Key`
- Mailgun: `MailGun_API_Key`

None of these `.env` files are committed — you'll need to create them yourself before starting the stack.
