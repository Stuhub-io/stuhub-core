FROM postgres:15-bullseye

COPY ./build/local/postgres/conf.d/postgresql.conf /etc/postgresql.conf
COPY ./build/local/postgres/initdb.d/ /docker-entrypoint-initdb.d/

# Maintainence scripts 
COPY ./build/local/postgres/maintenance /usr/local/bin/maintenance
RUN chmod +x /usr/local/bin/maintenance/*
# Make scripts available inin $path
RUN mv /usr/local/bin/maintenance/* /usr/local/bin \
    && rmdir /usr/local/bin/maintenance
