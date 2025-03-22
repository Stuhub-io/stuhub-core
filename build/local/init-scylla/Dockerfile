# Use the official ScyllaDB image
FROM scylladb/scylla

COPY ./build/scripts/wait-for-it.sh /wait-for-it.sh
RUN sed -i 's/\r$//g' /wait-for-it.sh
RUN chmod +x /wait-for-it.sh

# Copy scripts into the container
COPY ./build/local/init-scylla/entrypoint.sh /entrypoint.sh
RUN sed -i 's/\r$//g' /entrypoint.sh
RUN chmod +x /entrypoint.sh

COPY  ./build/local/init-scylla/start.sh /start.sh
RUN sed -i 's/\r$//g' /start.sh
RUN chmod +x /start.sh

ENTRYPOINT ["/entrypoint.sh"]
