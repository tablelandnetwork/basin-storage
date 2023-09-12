# basin-storage

[![License](https://img.shields.io/github/license/tablelandnetwork/basin-storage.svg)](./LICENSE)
[![standard-readme compliant](https://img.shields.io/badge/standard--readme-OK-green.svg)](https://github.com/RichardLitt/standard-readme)

> Archive data on Filecoin network.

# Table of Contents

- [basin-storage](#basin-storage)
- [Table of Contents](#table-of-contents)
- [Background](#background)
- [Development](#development)
  - [Running](#running)
  - [Run tests](#run-tests)  
- [Contributing](#contributing)
- [License](#license)

# Background

Tableland Basin is a secure and verifiable open data platform. The Basin Storage replicates data to Filecoin. It works in conjunction with [basin-cli](https://github.com/tablelandnetwork/basin-cli.git) and [basin-provider](https://github.com/tablelandnetwork/basin-provider.git).

🚧 Basin is currently not in a production-ready state. Any data that is pushed to the network may be subject to deletion. 🚧

# Development

Basin Storage leverages GCP Cloud Functions and Clould Storage to create a data archiving pipeline to Filecoin network. The File uploader function is triggered by an event from Cloud Storage.

This repository includes a development server that acts an a event trigger during development.

## Running

Start the development server for testing Clould Functions locally.
The required environment variables can be provided in `uploader.env.yml` and `checker.env.yml`.

```bash
make uploader-local
```

```bash
make checker-local
```

After the server is running, a mock cloud event can be triggered. It will execute the handler locally. In the following payload, we must have a real file path and a bucket for event simulation.

```bash
curl -X POST http://localhost:8080 \
-H "Content-Type: application/json" \
-H "ce-id: 1234567890" \
-H "ce-specversion: 1.0" \
-H "ce-type: google.cloud.storage.object.v1.finalized" \
-H "ce-time: 2020-08-08T00:11:44.895529672Z" \
-H "ce-source: //storage.googleapis.com/projects/_/buckets/tableland-entrypoint" \
-d '{
  "name": "feeds/2023-08-29/202308291525552525242120000000000-3ab461ed932d5f1c-1-2-00000000-employees-2.parquet",
  "bucket": "tableland-entrypoint",
  "contentType": "application/json",
  "metageneration": "1",
  "timeCreated": "2020-04-23T07:38:57.230Z",
  "updated": "2020-04-23T07:38:57.230Z"
}'
```

The checker function can be triggered by simply sending a POST request for example `curl -XPOST localhost:8080`.

## Deploying Function

```bash
make uploader-deploy
```

```bash
make checker-deploy
```

## Run tests

```bash
make test
```

# Contributing

PRs accepted.

Small note: If editing the README, please conform to the
[standard-readme](https://github.com/RichardLitt/standard-readme) specification.

# License

MIT AND Apache-2.0, © 2021-2023 Tableland Network Contributors
