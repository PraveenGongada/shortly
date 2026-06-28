# Shortly — Frontend

A minimal, dark, production-grade frontend for the Shortly URL shortener. Built
as a static React SPA (Vite + TypeScript + Tailwind) and served by nginx.

## Stack

- **Vite + React 18 + TypeScript** — static SPA, no SSR/Node server at runtime
- **Tailwind CSS** — dark, monochrome (Vercel-style) design; **Geist** font
- **TanStack Query** — link data fetching/caching/invalidation
- **react-hook-form + zod** — forms with validation mirroring the backend
- **Radix Dialog**, **sonner** (toasts), **lucide-react** (icons), **qrcode.react**
- **nginx** (unprivileged) — serves the build + SPA fallback

## Features

Mirrors the backend API exactly:

- Register / login / logout (JWT via HttpOnly cookie; register auto-logs-in)
- Create short links, list (paginated), edit destination, delete
- Per-link click count, copy, and a client-side **QR code**
- In-app short-link resolver at `/{code}` (forwards to the destination)

## Develop

```bash
yarn install
yarn dev             # http://localhost:3000, proxies /api -> $BACKEND_URL
```

Point the dev proxy at your backend (defaults to `http://localhost:8080`):

```bash
BACKEND_URL=http://localhost:8080 yarn dev
```

Copy `.env.example` to `.env` to set `VITE_BACKEND_URL` (leave empty for the
same-origin `/api` default).

## Build

```bash
yarn build           # type-checks then emits dist/
yarn preview         # serve the production build locally
```

## Docker

The image is a static build served by non-root nginx on **port 3000**.

```bash
docker build -t shortly-frontend .
docker run --rm -p 3000:3000 -e BACKEND_URL=http://host.docker.internal:8080 shortly-frontend
```

Build-time arg: `VITE_BACKEND_URL` (default empty → relative `/api`; set it to an
absolute URL like `https://api.shortly.com` to call that backend verbatim when
it's on a different host). Runtime env: `BACKEND_URL` for the optional `/api`
proxy (inert when an Istio gateway already routes `/api/*` to the backend).

## How requests work

The SPA calls a relative `/api`. With the existing Istio gateway routing
`/api/*` → backend and everything else → this app on one domain, the auth cookie
stays first-party and no CORS config is needed. Short links resolve through the
SPA (`/{code}` → `GET /api/{code}` → client redirect).
