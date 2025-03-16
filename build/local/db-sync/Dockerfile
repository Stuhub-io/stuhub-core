FROM postgres:15-bullseye

RUN apt-get update && apt-get install -y cron

WORKDIR /app

COPY ./build/local/db-sync/sync.sh /sync.sh
COPY ./build/local/db-sync/entrypoint.sh /entrypoint.sh

RUN chmod +x /sync.sh
RUN chmod +x /entrypoint.sh

ENTRYPOINT [ "/entrypoint.sh" ]
