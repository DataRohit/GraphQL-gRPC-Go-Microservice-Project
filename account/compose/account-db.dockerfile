FROM postgres:16

COPY ./account/sql/up.sql /docker-entrypoint-initdb.d/1.sql

CMD ["postgres"]
