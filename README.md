# Gator

Gator is a CLI tool made in Go to aggregate blog RSS feeds. 

# Install requirements

- Postgres v1.15+
- Go

# How to install

Clone the repo and run in your terminal (in the root directory of the project)

```bash
go install ...
```

Then create a `.gatorconfig.json` file in your home directory with the following structure:

```json
{
  "db_url": "postgres://username:@localhost:5432/database?sslmode=disable"
}
```

Replace the `db_url` value with your corresponding database connection string.

# Usage

## Users

### Register a new user

```bash
gator register <name>
```

### Login

```bash
gator login <name>
```

### List users

```bash
gator users
```

## Feeds

### Add a new feed

```bash
gator addfeed <name> <url>
```

### List feeds

```bash
gator feeds
```

### Follow a feed added by another user

```bash
gator follow <url>
```

### Unfollow a feed

```bash
gator unfollow <url>
```

### List user-feed follow relationships

```bash
gator following
```

## Aggregation

### Aggregate feeds

```bash
gator agg [frequency]
```

This command is meant to run in the background. It runs the aggregation process. It scrapes all the feeds that the currently logged in user follows and adds their posts to the database.

An optional argument can be passed to indicate the frequency in which to scrape each feed.

Example:

```bash
gator agg 1m
```
```bash
gator agg 30s
```

### Browse feeds

```bash
gator browse [limit]
```

Lists the posts that were aggregated with `gator agg` in chronological descending order.

An optional `limit` parameter can be added. If not, the default is 2.

Exmaple:

```bash
gator browse 5
```

## Misc

```bash
gator reset
```

Used for testing purposes only. It resets all the databases.

# Final words

This project was made for the "Build a Blog Aggregator" guided project from Boot.dev.