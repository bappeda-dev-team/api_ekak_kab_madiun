ARG GO_VERSION=1.22.8

FROM registry.docker.com/library/golang:$GO_VERSION-alpine as base

# app lives here
WORKDIR /app


# Throw-away build stage to reduce size of final image
FROM base as build

# Install packages needed to build
RUN apk update -qq && \
    apk add --no-cache git

COPY . .

RUN go build -o api main.go wire_gen.go

ENTRYPOINT ["/app/api"]

CMD ["app/api"]
