FROM redis:7-alpine

LABEL maintainer="Studio Senkou"
LABEL description="Production LC Redis"

RUN apk add --no-cache \
    tzdata \
    && rm -rf /var/cache/apk/*

ENV TZ=Asia/Jakarta
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

COPY redis.conf /usr/local/etc/redis/redis.conf

RUN mkdir -p /data && chown redis:redis /data

EXPOSE 6379

HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD redis-cli ping | grep PONG || exit 1

CMD ["redis-server", "/usr/local/etc/redis/redis.conf"]
