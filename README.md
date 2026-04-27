# The Coven — Community Platform

> Learn. Make. Share. Grow together.

The community platform for **[The Coven](https://thecoven.space)** — a collaborative makerspace for the queer community, founded by trans creators and focused on mutual support, skill-sharing, and creative projects in music, cosplay, video, and electronics.


## Repository Layout

| Path | Description |
| --- | --- |
| [hugo-site/](hugo-site/) | Public website built with [Hugo](https://gohugo.io) — tutorials, docs, project gallery, and landing pages. |
| [member-site/](member-site/) | The logged in website for members |
| [scripts/](scripts/) | Helper scripts. |

## The Hugo Site

The site lives in [hugo-site/](hugo-site/) and is deployed to **https://thecoven.space**.

### Prerequisites

- [Hugo (extended)](https://gohugo.io/installation/) — required for SCSS and image processing.

### Local Development

```bash
cd hugo-site
hugo server -D
```

The site will be available at `http://localhost:1313`. The `-D` flag includes draft content.

### Build

```bash
cd hugo-site
hugo --minify
```

The static output is written to [hugo-site/public/](hugo-site/public/).

### Content Structure

- [content/](hugo-site/content/) — Markdown content organized by section (about, charter, learning, projects, workshops, etc.).
- [data/](hugo-site/data/) — Structured YAML data backing many pages (FAQ, tour, code of conduct, etc.).
- [layouts/](hugo-site/layouts/) — Hugo templates and partials.
- [assets/](hugo-site/assets/) — CSS, JS, and images processed through Hugo Pipes.
- [static/](hugo-site/static/) — Files served as-is (CNAME, robots.txt, etc.).

## Roadmap

The current public site is the first piece of a larger system. See [plan.md](plan.md) for the full architecture, which includes:

- **Public site** (Hugo) — this repo
- **Member portal** (React SPA) — billing, prints, projects, access log
- **Backend API** (FastAPI or Go) — hardware control and Stripe integration
- **Database & Auth** — self-hosted Supabase (PostgreSQL + GoTrue + Storage + Realtime)
- **Edge & SSO** — Pangolin VPS reverse proxy
- **Physical access control** — Raspberry Pi + PN532 NFC + Dormakaba 8375 maglock
- **3D printer & camera monitoring** — PrusaLink polling and MJPEG/RTMP proxying

## Contributing

The Coven is a community space — contributions, corrections, and suggestions are welcome. Open an issue or pull request, or reach out via [hello@thecoven.space](mailto:hello@thecoven.space).

## License

[MIT](LICENSE) © atlantacoven
