# go-alexa-api

Voice skill server for controlling home entertainment devices (TVs, Roku, receiver) across multiple rooms. Supports both **Alexa** and **Google Assistant** via Dialogflow webhooks.

Built with Go using the [go-alexa](https://github.com/mikeflynn/go-alexa/tree/master/skillserver) skill server library.

## Prerequisites

- Go 1.25+
- Docker (for containerized deployment)

## Environment Variables

| Variable | Description |
|----------|-------------|
| `MBR_APP_ID` | Alexa App ID for master bedroom skill |
| `FR_APP_ID` | Alexa App ID for family room skill |
| `GOOGLE_WEBHOOK_TOKEN` | Shared secret for Google Assistant webhook auth (optional but recommended) |

Copy `.env.example` to `.env` and fill in your values.

## Build

```bash
go build -o go-alexa-api .
```

## Test

```bash
go test ./...
```

## Run Locally

```bash
export MBR_APP_ID=your-mbr-app-id
export FR_APP_ID=your-fr-app-id
./go-alexa-api
```

The server listens on port 8000.

## Docker

```bash
docker compose up -d
```

## Updating on the Production Server

The app is built from source inside Docker, so updating requires pulling the new code and rebuilding the image:

```bash
git pull
docker compose up -d --build
```

`--build` forces Docker to rebuild the image from the latest code before restarting the container. The old container is stopped and replaced automatically. Downtime is a few seconds during the rebuild.

To verify the new version is running:

```bash
docker compose ps
docker compose logs --tail=20
```

## Supported Voice Commands

| Command | Description |
|---------|-------------|
| OFF | Power off TV (and receiver in family room) |
| MUTE / UNMUTE | Toggle mute |
| VOLUME {level} | Set volume |
| CHANNEL {number} | Change channel |
| CHANNEL UP / DOWN | Channel up/down |
| INPUT {type} | Switch input (see below) |
| HOME / BACK | Roku navigation |
| UP / DOWN / LEFT / RIGHT {spaces} | Roku directional navigation |
| ENTER / SELECT | Roku confirm |
| PLAY / FORWARD / REVERSE | Roku playback |
| SEARCH {query} | Roku search |

## Supported Inputs

**Family Room:** TV, RetroPi, PS3, PS4, PS5, WiiU, FireTV/Roku, Switch, Xbox, Netflix, Plex, Prime, HBO, Crunchyroll, YouTube, and more.

**Master Bedroom:** TV, PS2, Wii, Switch, Netflix, Plex, Prime, HBO, Crunchyroll, YouTube, and more.

Each input supports multiple voice aliases (e.g., "Netflix", "Net", "Flix" all work).

---

## Google Assistant Setup

The server exposes Dialogflow webhook endpoints on the same port as Alexa:

| Endpoint | Room |
|----------|------|
| `POST /google/fr` | Family Room |
| `POST /google/mbr` | Master Bedroom |

### 1. Generate a Webhook Token

The `GOOGLE_WEBHOOK_TOKEN` is a secret you create — not issued by Google:

```bash
openssl rand -hex 32
```

Add the output to your `.env` file as `GOOGLE_WEBHOOK_TOKEN=<value>`.

### 2. Expose the Server via HTTPS

Dialogflow requires a public HTTPS webhook URL.

**For testing** — use [ngrok](https://ngrok.com):
```bash
ngrok http 8000
```

**For production** — put nginx or Caddy in front of this server with a TLS certificate (e.g., Let's Encrypt via Certbot).

### 3. Create a Dialogflow ES Agent

1. Go to [dialogflow.cloud.google.com](https://dialogflow.cloud.google.com) and create an agent.
2. Under **Fulfillment**, enable **Webhook** and set the URL:
   - Family Room: `https://your-domain.com/google/fr?token=<your-token>`
   - Master Bedroom: `https://your-domain.com/google/mbr?token=<your-token>`

> If you need both rooms, create two Dialogflow agents — one per room — each pointing to its own endpoint URL.

### 4. Create Intents

Create one intent per command. The **Intent name** must match the table below exactly. Enable **webhook fulfillment** on each intent.

| Intent Name | Parameters (Name → Entity) | Example Phrases |
|-------------|----------------------------|-----------------|
| `OFF` | — | "turn off the tv" |
| `MUTE` | — | "mute", "mute the tv" |
| `UNMUTE` | — | "unmute" |
| `VOLUME` | `Level` → `@sys.number` | "set volume to 30" |
| `CHANNEL` | `Number` → `@sys.number` | "go to channel 5" |
| `CHANNELUP` | — | "channel up" |
| `CHANNELDOWN` | — | "channel down" |
| `INPUT` | `InputType` → `@sys.any` | "switch to netflix" |
| `HOME` | — | "go home" |
| `BACK` | — | "go back" |
| `UP` | `Spaces` → `@sys.number` (optional) | "go up", "move up 3" |
| `DOWN` | `Spaces` → `@sys.number` (optional) | "go down" |
| `LEFT` | `Spaces` → `@sys.number` (optional) | "go left" |
| `RIGHT` | `Spaces` → `@sys.number` (optional) | "go right" |
| `ENTER` | — | "press enter" |
| `SELECT` | — | "select" |
| `PLAY` | — | "play" |
| `FORWARD` | — | "fast forward" |
| `REVERSE` | — | "rewind" |
| `SEARCH` | `SearchType` → `@sys.any` | "search for breaking bad" |

### 5. Connect to Google Assistant (optional)

In Dialogflow, go to **Integrations → Google Assistant** and click **Test**. This opens the Actions on Google simulator. For personal use you do not need to publish — keeping the action in draft/test mode works on your own Google account's devices indefinitely.
