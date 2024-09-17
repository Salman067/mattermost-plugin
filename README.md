# Mattermost Docker Deployment
The Docker deployment solution for Mattermost.

## Install & Usage

Refer to the [Mattermost Docker deployment guide](https://docs.mattermost.com/install/install-docker.html) for instructions on how to install and use this Docker image.

## Upgrading from `mattermost-docker`

This repository replaces the [deprecated mattermost-docker repository](https://github.com/mattermost/mattermost-docker). For an in-depth guide to upgrading, please refer to [this document](https://github.com/mattermost/docker/blob/main/scripts/UPGRADE.md).

## TO run locally

```
docker compose -f docker-compose.yml -f docker-compose.without-nginx.yml up -d  
```
