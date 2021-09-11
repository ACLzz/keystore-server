FROM golang:1.17
RUN mkdir /app
ADD . /app
COPY .docker_env /app/bin/.env
COPY .docker_env /app/.env
WORKDIR /app
RUN apt-get update \
 && DEBIAN_FRONTEND=noninteractive \
    apt-get install --no-install-recommends -yq \
      postgresql-client
RUN make setup
EXPOSE 8402
CMD ["/app/bin/server"]