FROM docker.io/golang:alpine AS build
COPY . /tmp/src
WORKDIR /tmp/src
RUN mkdir -p /tmp/build
RUN go mod download -x
RUN go build -x -tags=docker,nomsgpack,go_json -ldflags "-w -s" -o /tmp/build/app

FROM docker.io/alpine:latest

ARG GH_REPO=unset
ARG GH_VERSION=unset
LABEL org.opencontainers.image.source=https://github.com/$GH_REPO
LABEL org.opencontainers.image.version=$GH_VERSION

COPY --from=build /tmp/build/app /service
ENTRYPOINT ["/service"]
HEALTHCHECK --interval=30s --timeout=15s CMD /service -healthcheck
EXPOSE 8000