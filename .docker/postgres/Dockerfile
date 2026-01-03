FROM postgres:17-alpine

LABEL maintainer="Studio Senkou"
LABEL description="Production LC PostgreSQL image"

RUN apk add --no-cache \
    tzdata \
    && rm -rf /var/cache/apk/*

ENV TZ=Asia/Jakarta
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

ENV LANG=en_US.UTF-8
ENV LC_ALL=en_US.UTF-8

COPY postgresql.conf /etc/postgresql/postgresql.conf
COPY pg_hba.conf /etc/postgresql/pg_hba.conf

COPY init/ /docker-entrypoint-initdb.d/

RUN chmod +x /docker-entrypoint-initdb.d/*.sh 2>/dev/null || true

RUN mkdir -p /var/lib/postgresql/data
CMD ["postgres", "-c", "config_file=/etc/postgresql/postgresql.conf", "-c", "hba_file=/etc/postgresql/pg_hba.conf"]
