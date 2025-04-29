FROM postgres:17 AS pg
LABEL maintainer="kaffein"

WORKDIR /root

RUN apt-get update && apt-get install -y wget
RUN wget https://github.com/pksunkara/pgx_ulid/releases/download/v0.2.0/pgx_ulid-v0.2.0-pg17-amd64-linux-gnu.deb
RUN dpkg -i pgx_ulid-v0.2.0-pg17-amd64-linux-gnu.deb
CMD ["postgres"]
