# Go Redis Facade

[![Validate](https://github.com/coopnorge/go-redis-facade/actions/workflows/cicd.yaml/badge.svg)](https://github.com/coopnorge/go-redis-facade/actions/workflows/cicd.yaml)

Coop Redis Facade wraps simple interaction with Redis clients for CRUD
operations by preventing race conditions between multiple client instances
against singular instances of Redis.

If you are interested in how Sync between clients works, take a look at [this
post.](https://redis.io/docs/manual/patterns/distributed-locks/)

## Module Documentation

<https://pkg.go.dev/github.com/coopnorge/go-redis-facade>

## Mocks

To generate or update mocks use
[`gomockhandler`](github.com/sanposhiho/gomockhandler). `gomockhandler` is
provided by `golang-devtools`.

### Check mocks

```bash
docker compose run --rm golang-devtools gomockhandler -config ./gomockhandler.json check
```

### Generate / Update mocks

```bash
docker compose run --rm golang-devtools gomockhandler -config ./gomockhandler.json mockgen
```
