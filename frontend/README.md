# Megan Messenger frontend

Vue 3 client built with FSD. The UI reads domain data only from Dexie/IndexedDB; TanStack Query standardizes HTTP state and writes successful payloads into the local database. Pinia holds only ephemeral UI/session state.

## Run

```bash
npm install
npm run dev
```

Vite proxies `/api/*` to `http://localhost:8080/*`.

Auth: http://localhost:5173/auth — phone → OTP (код в логах API) → username.

## Layout

- `src/app` — bootstrap, router, global styles
- `src/pages` — route screens (`auth`, `messenger`)
- `src/widgets` — composed UI
- `src/features` — user actions (`auth` flow, conversation selection)
- `src/entities` — domain (`session`, `conversation`) + API/Dexie
- `src/shared` — HTTP, Query client, DB, UI primitives

Keep server synchronization in `entities/*/api/*.queries.ts`; Vue components subscribe with `useLiveQuery` only.

## Vue SFC convention

Keep Vue single-file components in this order: `<template>`, then `<script setup lang="ts">`, then an optional `<style>`. This order applies to every new or modified `.vue` file.
