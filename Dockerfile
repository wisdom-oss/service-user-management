FROM docker.io/library/golang:alpine AS build
COPY . .
RUN --mount=type=cache,target=/go/pkg/mod go mod download -x
RUN --mount=type=cache,target=/go/pkg/mod --mount=type=cache,target=/root/.cache/go-build go build -x -tags=docker,nomsgpack,go_json -ldflags "-w -s" -o /service

FROM docker.io/library/alpine:latest

ARG GH_REPO=unset
ARG GH_VERSION=unset
LABEL org.opencontainers.image.source=https://github.com/$GH_REPO
LABEL org.opencontainers.image.version=$GH_VERSION

COPY --from=build /service /service
ENTRYPOINT ["/service"]
HEALTHCHECK --interval=30s --timeout=15s CMD /service -healthcheck
EXPOSE 8000