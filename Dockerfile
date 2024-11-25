FROM docker.io/golang:alpine AS build-service
ENV GOMODCACHE=/root/.cache/go-build
WORKDIR /src
COPY --link go.* .
RUN --mount=type=cache,target=/root/.cache/go-build go mod download
COPY --link . .
RUN --mount=type=cache,target=/root/.cache/go-build go build -tags=docker,nomsgpack,go_json -o /service .

FROM docker.io/alpine:latest

ARG GH_REPO=unset
ARG GH_VERSION=unset
LABEL org.opencontainers.image.source=https://github.com/$GH_REPO
LABEL org.opencontainers.image.version=$GH_VERSION

COPY --link --from=build-service /service /service
ENTRYPOINT ["/service"]
HEALTHCHECK --interval=30s --timeout=15s CMD /service -healthcheck
EXPOSE 8000

