FROM edoburu/pgbouncer:latest

COPY pgbouncer.ini /etc/pgbouncer/pgbouncer.ini
COPY --chmod=755 entrypoint.sh ./entrypoint.sh

ENTRYPOINT ["./entrypoint.sh"]

