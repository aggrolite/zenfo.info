# zenfo.info

zenfo.info crawls Zen center websites for event information.

Events are stored into a Postgres database, and served as JSON via HTTP.

## Sources

So far there are two sources:

1. San Francisco Zen Center (sfzc.org)
2. Angel City Zen Center (aczc.org)

Sources are scraped by using the `Worker` interface. All workers share a common HTTP client which provides a custom user agent. And in the future there will need to be rate limiting per site.

## Milestones for 1.0 release

1. Add at least two more sources for crawling. Ideally one from East Bay, and one South Bay.
2. Add basic frontend UI, ideally something in JS.
3. Add some sort of detection on when a crawler is completely busted and needs updating.
4. Add Dockerfile for easy deploy.

## Build / Install

For now, simply run `make` to try it out. A Postgres instance is required to be running in the background using provided .psql schema.

## License

Currently this project is licensed under AGPL v3.0. See `LICENSE` for details.

I'm usually not a fan of anything GPL, but with this project it feels right. License may change in the future to something more permissive.
