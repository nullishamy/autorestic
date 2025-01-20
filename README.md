# autorestic

Run your restic backups with logging and error handling

## usage
1) Set the required environment variables
  - `RESTIC_REPOSITORY`    (your restic repository)
  - `RESTIC_PASSWORD`      (your restic password)
  - `AUTORESTIC_WEBHOOK`   (your Dicord webhook)
  - `AUTORESTIC_LOCATIONS` (the directories to backup from)
  - `AUTORESTIC_IGNORE`    (the directories to ignore [optional])

2) Run autorestic on a schedule (e.g crontab)
3) Profit!

## building
You will probably want to build this as a static binary:
- `CGO_ENABLED=0 go build  -ldflags="-extldflags=-static" *.go`


