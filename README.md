# planshare

Uses [pev2](https://github.com/dalibo/pev2) and adds simple go backend to allow sharing plan visualisations.

Requires redis to store data.

## How to run on local machine

```
docker run -p 8080:80 -e REDIS_ADDR=172.17.0.1:6379 mluvii/planshare
```
