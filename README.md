# go-alexa-api

Alexa skill server for controlling home entertainment devices (TVs, Roku, receiver) across multiple rooms via voice commands.

Built with Go using the [go-alexa](https://github.com/mikeflynn/go-alexa/tree/master/skillserver) skill server library.

## Prerequisites

- Go 1.25+
- Docker (for containerized deployment)

## Environment Variables

| Variable | Description |
|----------|-------------|
| `MBR_APP_ID` | Alexa App ID for master bedroom skill |
| `FR_APP_ID` | Alexa App ID for family room skill |

Copy `.env.example` to `.env` and fill in your app IDs.

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

## Deploy to Raspberry Pi

```bash
./deploy.sh
```

Builds an ARM image, transfers it to the media server, and starts the container.

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
