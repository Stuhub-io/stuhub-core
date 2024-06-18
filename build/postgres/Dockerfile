FROM postgres:16

COPY build/postgres/conf.d/postgresql.conf /etc/postgresql.conf
COPY build/postgres/initdb.d/ /docker-entrypoint-initdb.d/