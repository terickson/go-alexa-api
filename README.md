# go-alexa-api

Voice skill server for controlling home entertainment devices (TVs, Roku, receiver) across multiple rooms. Supports **Alexa** and **Google Assistant via IFTTT**.

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

## Google Assistant via IFTTT Setup

> **Note:** Google Assistant Conversational Actions was permanently shut down on **June 13, 2023**. The Dialogflow → Google Assistant integration path no longer exists. The replacement is **IFTTT** — it still has a working Google Assistant trigger and can call any URL via its Webhooks action.

The server exposes simple GET endpoints designed for IFTTT. All parameters go in the URL as query params — no JSON body required.

| Endpoint | Room |
|----------|------|
| `GET /action/fr` | Family Room |
| `GET /action/mbr` | Master Bedroom |

### 1. Generate a Webhook Token

```bash
openssl rand -hex 32
```

Add the output to your `.env` file as `GOOGLE_WEBHOOK_TOKEN=<value>`. If this variable is set, every request must include `?token=<value>` or it will be rejected with 401.

### 2. Expose the Server via HTTPS

IFTTT requires a public HTTPS URL.

**For production** — put nginx or Caddy in front of this server with a TLS certificate (e.g., Let's Encrypt via Certbot).

### 3. Create an IFTTT Account

Go to [ifttt.com](https://ifttt.com) and sign in with your Google account. **IFTTT Pro** (~$2.99/month) is required if you need more than 2 applets.

### 4. Create an Applet

For each voice command, create one applet:

1. **If This** → Search for **Google Assistant** → Choose a trigger:
   - "Say a simple phrase" — for commands with no variable (OFF, MUTE, HOME, etc.)
   - "Say a phrase with a number ingredient" — for VOLUME, CHANNEL, spaces
   - "Say a phrase with a text ingredient" — for INPUT, SEARCH
2. Fill in the trigger phrase(s). The ingredient placeholder is `{{NumberField}}` or `{{TextField}}`.
3. **Then That** → Search for **Webhooks** → **Make a web request**
   - Method: **GET**
   - URL: see examples below
   - Content Type: `application/json`
   - Body: (leave empty)

### 5. Example Applet URLs

Replace `YOUR_DOMAIN` with your server's hostname and `YOUR_TOKEN` with your `GOOGLE_WEBHOOK_TOKEN` value.

#### Family Room

| Voice phrase | URL |
|---|---|
| "turn off the TV in the family room" | `https://YOUR_DOMAIN/action/fr?token=YOUR_TOKEN&intent=OFF` |
| "mute the TV in the family room" | `https://YOUR_DOMAIN/action/fr?token=YOUR_TOKEN&intent=MUTE` |
| "unmute the TV in the family room" | `https://YOUR_DOMAIN/action/fr?token=YOUR_TOKEN&intent=UNMUTE` |
| "set family room volume to #" (NumberField) | `https://YOUR_DOMAIN/action/fr?token=YOUR_TOKEN&intent=VOLUME&level={{NumberField}}` |
| "go to channel # in the family room" (NumberField) | `https://YOUR_DOMAIN/action/fr?token=YOUR_TOKEN&intent=CHANNEL&number={{NumberField}}` |
| "channel up in the family room" | `https://YOUR_DOMAIN/action/fr?token=YOUR_TOKEN&intent=CHANNELUP` |
| "channel down in the family room" | `https://YOUR_DOMAIN/action/fr?token=YOUR_TOKEN&intent=CHANNELDOWN` |
| "switch the family room to $" (TextField) | `https://YOUR_DOMAIN/action/fr?token=YOUR_TOKEN&intent=INPUT&input={{TextField}}` |
| "search the family room for $" (TextField) | `https://YOUR_DOMAIN/action/fr?token=YOUR_TOKEN&intent=SEARCH&search={{TextField}}` |
| "go home in the family room" | `https://YOUR_DOMAIN/action/fr?token=YOUR_TOKEN&intent=HOME` |
| "go back in the family room" | `https://YOUR_DOMAIN/action/fr?token=YOUR_TOKEN&intent=BACK` |
| "go up in the family room" | `https://YOUR_DOMAIN/action/fr?token=YOUR_TOKEN&intent=UP` |

#### Master Bedroom

Same as above with `/action/mbr` instead of `/action/fr`.

### Intent Reference

| Intent | Extra Query Param | Default |
|--------|------------------|---------|
| `OFF`, `MUTE`, `UNMUTE`, `HOME`, `BACK`, `ENTER`, `SELECT`, `PLAY`, `FORWARD`, `REVERSE`, `CHANNELUP`, `CHANNELDOWN` | none | — |
| `VOLUME` | `level=<number>` | required |
| `CHANNEL` | `number=<number>` | required |
| `INPUT` | `input=<text>` | required |
| `SEARCH` | `search=<text>` | required |
| `UP`, `DOWN`, `LEFT`, `RIGHT` | `spaces=<number>` | `1` |

### Testing

You can test any command by pasting a URL directly in your browser (or with curl):

```bash
curl "https://YOUR_DOMAIN/action/fr?token=YOUR_TOKEN&intent=OFF"
# → {"status":"ok","intent":"OFF"}

curl "https://YOUR_DOMAIN/action/fr?token=YOUR_TOKEN&intent=VOLUME&level=25"
# → {"status":"ok","intent":"VOLUME"}

curl "https://YOUR_DOMAIN/action/fr?token=YOUR_TOKEN&intent=INPUT&input=netflix"
# → {"status":"ok","intent":"INPUT"}
```
