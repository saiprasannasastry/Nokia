FROM postgres 
ENV POSTGRES_PASSWORD postgres 
ENV POSTGRES_DB album
ENV POSTGRES_USER postgres
COPY psql/init.sql /docker-entrypoint-initdb.d/
