# The Hugo Site

The site lives in [hugo-site/](hugo-site/) and is deployed to **https://atlantacoven.org**.

## Prerequisites

- [Hugo (extended)](https://gohugo.io/installation/) — required for SCSS and image processing.

## Local Development

```bash
cd hugo-site
hugo server -D
```

The site will be available at `http://localhost:1313`. The `-D` flag includes draft content.

## Build

```bash
cd hugo-site
hugo --minify
```

The static output is written to [hugo-site/public/](hugo-site/public/).

## Content Structure

- [content/](hugo-site/content/) — Markdown content organized by section (about, charter, learning, projects, workshops, etc.).
- [data/](hugo-site/data/) — Structured YAML data backing many pages (FAQ, tour, code of conduct, etc.).
- [layouts/](hugo-site/layouts/) — Hugo templates and partials.
- [assets/](hugo-site/assets/) — CSS, JS, and images processed through Hugo Pipes.
- [static/](hugo-site/static/) — Files served as-is (CNAME, robots.txt, etc.).
