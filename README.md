# dlog

`dlog` is a small CLI for recording and viewing daily work logs.

## Install

```bash
go install github.com/okaryo/dlog@latest
```

## Usage

```bash
dlog "task progress update"
dlog
dlog log
dlog log --date 2026-04-12
dlog md
dlog md --date 2026-04-12
```

## Behavior

- `dlog "text"` appends a log entry for today with the current local timestamp.
- `dlog` and `dlog log` show today's logs in reverse chronological order.
- `dlog log --date YYYY-MM-DD` shows logs for the specified date.
- `dlog md` prints today's logs in Markdown order from oldest to newest.
- `dlog md --date YYYY-MM-DD` prints logs for the specified date in Markdown order.

## Storage

- Logs are stored under `~/.dlog`.
- Each day is stored as one JSON file named `YYYY-MM-DD.json`.
- Timestamps are saved in RFC3339 format with the local time zone.

Example:

```json
{
  "date": "2026-04-12",
  "logs": [
    {
      "timestamp": "2026-04-12T10:03:21+09:00",
      "text": "task progress update"
    }
  ]
}
```

## Future ideas

- Additional export formats
- Daily and weekly summaries
- Tags
