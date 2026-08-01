# Megan Messenger frontend

Vue 3 client built with FSD. The UI reads domain data only from Dexie/IndexedDB; TanStack Query standardizes HTTP state and writes successful payloads into the local database. Pinia holds only ephemeral UI/session state.

## Run

```bash
npm install
npm run dev
```

Vite proxies `/api/*` to `http://localhost:8080/*`. Set `VITE_API_URL` later if a deployed API needs a different transport strategy.

## Layout

- `src/app` — application bootstrap and global styles
- `src/entities` — domain types, API adapters, and Dexie persistence
- `src/features` — user actions and Pinia UI state
- `src/widgets` — composed UI
- `src/shared` — API transport, Query client, DB, reusable utilities

Keep server synchronization in `entities/*/api/*.queries.ts`; Vue components subscribe with `useLiveQuery` only.

## Vue SFC convention

Keep Vue single-file components in this order: `<template>`, then `<script setup lang="ts">`, then an optional `<style>`. This order applies to every new or modified `.vue` file.
