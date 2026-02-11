# Gator — CLI RSS aggregator

> Lightweight command-line RSS aggregator and personal feed manager.

This project provides a small CLI to manage users, follow RSS feeds, run a background aggregator, and browse cached posts.

## What it does

- Manage local users (`register`, `login`, `users`, `reset`).
- Add and follow RSS feeds (`addfeed`, `follow`, `following`, `unfollow`, `feeds`).
- Aggregate feeds periodically and store posts in a Postgres database (`agg`).
- Browse stored posts (`browse`).

## Prerequisites

- Go (1.18+ recommended)
- PostgreSQL database

Make sure you have the repository tracked with Git (e.g. `git init` and commit), so you can manage changes.

## Quick setup

1. Create a Postgres database:

```bash
createdb gator
```

2. Apply the SQL schema (simple approach using `psql`):

```bash
for f in sql/schema/*.sql; do psql -d gator -f "$f"; done
```

3. Create a config file at `~/.gatorconfig.json` with your DB URL. Example:

```json
{
  "dbUrl": "postgres://user:pass@localhost:5432/gator?sslmode=disable",
  "current_user_name": ""
}
```

Alternatively, you can create a user via the CLI which will update this file automatically.

## Build

```bash
go build -o gator .
```

Run the binary from the repository root (or install it to your PATH).

## Install via `go install`

You can install the `gator` CLI with `go install` from anywhere (requires Go modules and your `GOBIN`/`PATH` configured):

```bash
go install github.com/j-machuca/gator@latest
```

After installing, run the command with `gator` (or the binary name you built).

## Config file

The CLI reads/writes a small config file at `~/.gatorconfig.json`. The config should contain at least your database URL. Example file contents:

```json
{
  "dbUrl": "postgres://user:pass@localhost:5432/gator?sslmode=disable",
  "current_user_name": ""
}
```

You can create that file manually, or the CLI will update `current_user_name` when you run `register` or `login`.

## Running the program

Run the binary or installed `gator` with a command, for example:

```bash
./gator register alice
./gator login alice
./gator addfeed "Go Blog" "https://blog.golang.org/feed.atom"
./gator agg 30s   # starts the aggregator loop (runs until interrupted)
./gator browse 5  # show up to 5 recent posts
```

If you installed with `go install`, replace `./gator` with `gator`.

## Common commands (quick reference)

- `register <name>` — create a new user and set it as current.
- `login <name>` — switch current user.
- `addfeed <name> <url>` — create a feed and follow it (must be logged in).
- `agg <duration>` — run continuous aggregation loop fetching feeds every `<duration>` (e.g. `10s`, `1m`).
- `browse [limit]` — show the most recent stored posts (default limit: 2).

## Usage

General form:

```bash
./gator <command> [args...]
```

Available commands (registered in `main.go`):

- `register <name>` — create a new user and set it as current.
- `login <name>` — switch current user (user must exist).
- `reset` — clear all users, feeds and followings (and unset current user).
- `users` — list users (current user is marked).
- `addfeed <name> <url>` — create a feed and follow it (requires logged in user).
- `feeds` — list all feeds (name, url, owner username).
- `follow <url>` — follow an existing feed by URL (requires logged in user).
- `following` — list feeds the current user follows.
- `unfollow <url>` — unfollow a feed by URL (requires logged in user).
- `agg <duration>` — run an infinite aggregation loop fetching feeds every `<duration>` (e.g. `10s`).
- `browse [limit]` — show the most recent stored posts (default limit: 2).

## Examples and sample outputs

Note: some commands interact with the database and depend on prior state. The following examples show representative outputs.

- Register a user:

```bash
$ ./gator register alice
user succesfully created user:
 * ID:      0a1b2c3d-...-abcd
 * Name:    alice
```

- Login as an existing user:

```bash
$ ./gator login alice
User switched succesfully!
```

- List users (current marked):

```bash
$ ./gator users
* alice (current)
* bob
```

- Reset the database:

```bash
$ ./gator reset
reset was successful!
```

- Add a feed (must be logged in):

```bash
$ ./gator addfeed "Go Blog" "https://blog.golang.org/feed.atom"
Feed:
 * ID:01234567-...-89ab
 * Name:Go Blog
 * URL:https://blog.golang.org/feed.atom
 * UserId:ef012345-...-abcd
```

- List feeds:

```bash
$ ./gator feeds
Go Blog
https://blog.golang.org/feed.atom
alice
```

- Follow a feed (must be logged in):

```bash
$ ./gator follow https://blog.golang.org/feed.atom
Feed Name Go Blog
Current Username alice
```

- List followings for current user:

```bash
$ ./gator following
Go Blog
```

- Unfollow a feed:

```bash
$ ./gator unfollow https://blog.golang.org/feed.atom
# (no output on success)
```

- Run the aggregator (example `10s`). This starts an infinite loop; sample console output:

```bash
$ ./gator agg 10s
Collecting feeds every 10s
Starting the scrape loop
Fetching next feed to fetch
Marking https://blog.golang.org/feed.atom as fetched
Fetching data from https://blog.golang.org/feed.atom
iterating over items in feed
Inserted Post:
Title: An interesting article
Description: Summary of the article
URL: https://blog.golang.org/some-article
```

- Browse stored posts (default 2):

```bash
$ ./gator browse
Title: An interesting article
Description: Summary of the article
Published Date: Sat, 01 Jan 2000 00:00:00 +0000
URL: https://blog.golang.org/some-article
```

## Notes

- The CLI stores current username in `~/.gatorconfig.json` (see `internal/config`). Some commands require a logged-in user.
- Aggregation (`agg`) will continuously loop and insert posts into the `posts` table; stop it with Ctrl+C.
- Migrations are provided under `sql/schema` — they can be applied in order. The repository expects a Postgres backend and uses code generated by `sqlc` under `internal/database`.
