# PostgreSQL Backup
PostgreSQL Backup is a Docker container image that can be used to backup, restore and migrate Postgres database. It supports local storage, AWS S3 or any S3 Alternatives for Object Storage, FTP and SSH compatible storage.
It also supports __encrypting__ your backups using GPG.

The [jkaninda/pg-bkup](https://hub.docker.com/r/jkaninda/pg-bkup) Docker image can be deployed on Docker, Docker Swarm and Kubernetes.
It handles __recurring__ backups of postgres database on Docker and can be deployed as __Job or CronJob  on Kubernetes__ using local, AWS S3 or SSH compatible storage.

It also supports __encrypting__ your backups using GPG.

Telegram and Email notifications on successful and failed backups.


[![Build](https://github.com/jkaninda/pg-bkup/actions/workflows/release.yml/badge.svg)](https://github.com/jkaninda/pg-bkup/actions/workflows/release.yml)
[![Go Report](https://goreportcard.com/badge/github.com/jkaninda/mysql-bkup)](https://goreportcard.com/report/github.com/jkaninda/pg-bkup)
![Docker Image Size (latest by date)](https://img.shields.io/docker/image-size/jkaninda/pg-bkup?style=flat-square)
![Docker Pulls](https://img.shields.io/docker/pulls/jkaninda/pg-bkup?style=flat-square)
<a href="https://ko-fi.com/jkaninda"><img src="https://uploads-ssl.webflow.com/5c14e387dab576fe667689cf/5cbed8a4ae2b88347c06c923_BuyMeACoffee_blue.png" height="20" alt="buy ma a coffee"></a>


Successfully tested on:
- Docker
- Docker in Swarm mode
- Kubernetes
- OpenShift

## Documentation is found at <https://jkaninda.github.io/pg-bkup>


## Links:

- [Docker Hub](https://hub.docker.com/r/jkaninda/pg-bkup)
- [Github](https://github.com/jkaninda/pg-bkup)

## MySQL solution :

- [MySQL](https://github.com/jkaninda/mysql-bkup)

## Storage:
- Local
- AWS S3 or any S3 Alternatives for Object Storage
- SSH remote storage server
- FTP remote storage server

## Quickstart

### Simple backup using Docker CLI

To run a one time backup, bind your local volume to `/backup` in the container and run the `backup` command:

```shell
 docker run --rm --network your_network_name \
 -v $PWD/backup:/backup/ \
 -e "DB_HOST=dbhost" \
 -e "DB_USERNAME=username" \
 -e "DB_PASSWORD=password" \
 jkaninda/pg-bkup backup -d database_name
```

Alternatively, pass a `--env-file` in order to use a full config as described below.

```yaml
 docker run --rm --network your_network_name \
 --env-file your-env-file \
 -v $PWD/backup:/backup/ \
 jkaninda/pg-bkup backup -d database_name
```

### Simple backup in docker compose file

```yaml
services:
  pg-bkup:
    # In production, it is advised to lock your image tag to a proper
    # release version instead of using `latest`.
    # Check https://github.com/jkaninda/pg-bkup/releases
    # for a list of available releases.
    image: jkaninda/pg-bkup
    container_name: pg-bkup
    command: backup
    volumes:
      - ./backup:/backup
    environment:
      - DB_PORT=5432
      - DB_HOST=postgres
      - DB_NAME=foo
      - DB_USERNAME=bar
      - DB_PASSWORD=password
      - TZ=Europe/Paris
    # pg-bkup container must be connected to the same network with your database
    networks:
       - web
networks:
  web:
```
### Docker recurring backup

```shell
 docker run --rm --network network_name \
 -v $PWD/backup:/backup/ \
 -e "DB_HOST=hostname" \
 -e "DB_USERNAME=user" \
 -e "DB_PASSWORD=password" \
 jkaninda/pg-bkup backup -d dbName --cron-expression "@every 15m" #@midnight
```
See: https://jkaninda.github.io/pg-bkup/reference/#predefined-schedules

## Deploy on Kubernetes

For Kubernetes, you don't need to run it in scheduled mode. You can deploy it as Job or CronJob.

### Simple Kubernetes backup Job :

```yaml
apiVersion: batch/v1
kind: Job
metadata:
  name: backup-job
spec:
  ttlSecondsAfterFinished: 100
  template:
    spec:
      containers:
        - name: pg-bkup
          # In production, it is advised to lock your image tag to a proper
          # release version instead of using `latest`.
          # Check https://github.com/jkaninda/pg-bkup/releases
          # for a list of available releases.
          image: jkaninda/pg-bkup
          command:
            - /bin/sh
            - -c
            - backup -d dbname
          resources:
            limits:
              memory: "128Mi"
              cpu: "500m"
          env:
            - name: DB_HOST
              value: "postgres"
            - name: DB_USERNAME
              value: "postgres"
            - name: DB_PASSWORD
              value: "password"
          volumeMounts:
            - mountPath: /backup
              name: backup
      volumes:
        - name: backup
          hostPath:
            path: /home/toto/backup # directory location on host
            type: Directory # this field is optional
      restartPolicy: Never
```
## Available image registries

This Docker image is published to both Docker Hub and the GitHub container registry.
Depending on your preferences and needs, you can reference both `jkaninda/pg-bkup` as well as `ghcr.io/jkaninda/pg-bkup`:

```
docker pull jkaninda/pg-bkup
docker pull ghcr.io/jkaninda/pg-bkup
```

Documentation references Docker Hub, but all examples will work using ghcr.io just as well.

## References

We decided to publish this image as a simpler and more lightweight alternative because of the following requirements:

- The original image is based on `Alpine` and requires additional tools, making it heavy.
- This image is written in Go.
- `arm64` and `arm/v7` architectures are supported.
- Docker in Swarm mode is supported.
- Kubernetes is supported.


## License

This project is licensed under the MIT License. See the LICENSE file for details.

## Authors

**Jonas Kaninda**
- <https://github.com/jkaninda>

## Copyright

Copyright (c) [2023] [Jonas Kaninda]
