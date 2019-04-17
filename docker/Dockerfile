FROM library/postgres:9.6
COPY init.sql /docker-entrypoint-initdb.d/

ENV POSTGRES_USER navpool
ENV POSTGRES_PASSWORD navpool
ENV POSTGRES_DB navpool

EXPOSE 5432
