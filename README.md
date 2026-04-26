# dlog

`dlog` is a small CLI for daily work log.

## Install

```bash
go install github.com/okaryo/dlog@latest
```

## Usage

```bash
dlog "task progress update"
dlog amend "corrected task progress update"
dlog
dlog log
dlog log yesterday
dlog log --date 2026-04-12
dlog md
dlog md yesterday
dlog md --date 2026-04-12
```

## Behavior

- `dlog "text"` appends a log entry for today with the current local timestamp.
- `dlog amend "text"` replaces today's most recent log entry while keeping its original timestamp.
- `dlog` and `dlog log` show today's logs in reverse chronological order.
- `dlog log [YYYY-MM-DD|today|yesterday]` and `dlog log --date [YYYY-MM-DD|today|yesterday]` show logs for the specified date.
- `dlog md` prints today's logs in Markdown order from oldest to newest.
- `dlog md [YYYY-MM-DD|today|yesterday]` and `dlog md --date [YYYY-MM-DD|today|yesterday]` print logs for the specified date in Markdown order.
- Displayed times use the timezone recorded in each log entry, not the viewer's current local timezone.

## Examples

Record logs for today:

```bash
$ dlog "task progress update"
$ dlog "api design"
$ dlog "bug fix"
```

Amend the most recent log entry:

```bash
$ dlog amend "bug fix (retry with test case)"
```

Show today's logs in reverse chronological order with `dlog` or `dlog log`:

```bash
$ dlog
2026-04-12

11:10 bug fix (retry with test case)
10:25 api design
10:03 task progress update
```

```bash
$ dlog log
2026-04-12

11:10 bug fix (retry with test case)
10:25 api design
10:03 task progress update
```

Show a specific date:

```bash
$ dlog log --date 2026-04-11
2026-04-11

18:45 previous day task
```

Output logs as Markdown in chronological order:

```bash
$ dlog md
# 2026-04-12
- 10:03 task progress update
- 10:25 api design
- 11:10 bug fix (retry with test case)
```

Output a specific date as Markdown:

```bash
$ dlog md --date 2026-04-11
# 2026-04-11
- 18:45 previous day task
```

Output yesterday's logs as Markdown:

```bash
$ dlog md yesterday
# 2026-04-11
- 18:45 previous day task
```

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
