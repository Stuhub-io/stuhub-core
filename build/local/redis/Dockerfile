FROM redis:7.2-alpine

WORKDIR /redis

COPY ./build/local/redis/init.sh ./

RUN sed -i 's/\r$//g' init.sh
RUN chmod +x init.sh

