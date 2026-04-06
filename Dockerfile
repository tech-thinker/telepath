# Builder stage
FROM golang:1.25 AS builder

ARG VERSION="v0.0.0"
ARG COMMIT="Unknown"
ARG BUILDDATE="Unknown"

WORKDIR /src
ADD . .
RUN go build -gcflags="all=-N -l" -ldflags="-X 'main.AppVersion=${VERSION}' -X 'main.CommitHash=${COMMIT}' -X 'main.BuildDate=${BUILDDATE}'" -o telepath

# Runner Stage
FROM debian:bookworm-slim
# Install SQLite and required shared libs
RUN apt-get update && apt-get install -y libsqlite3-0 tzdata iputils-ping && \
	rm -rf /var/lib/apt/lists/*
RUN mkdir -p /app /etc/telepath
WORKDIR /app
COPY --from=builder /src/telepath /app/telepath
RUN chmod +x /app/telepath

ENTRYPOINT ["/app/telepath"]
CMD ["--config-file=/etc/telepath/config.json"]
