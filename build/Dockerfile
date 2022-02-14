FROM golang:1.15-alpine AS build
RUN apk update && \
    apk add curl \
    git \
    bash \
    make \
    ca-certificates && \
    rm -rf /var/cache/apk/*

RUN apk add --no-cache tzdata

WORKDIR /app

# copy module files first so that they don't need to be downloaded again if no change
ENV GOPATH /go
COPY go.* ./
RUN go mod download
RUN go mod verify


# copy source files and build the binary
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o go-clean-arch main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates bash

# RUN mkdir -p /var/log/app
WORKDIR /app/
COPY --from=build /app/go-clean-arch .
COPY --from=build /app/scripts/migrations/postgres ./scripts/migrations/postgres/
# COPY --from=build /app/assets ./assets/
RUN ls -la
ENV TZ Asia/Jakarta
RUN ["chmod", "+x", "./go-clean-arch"]
ENTRYPOINT ["./go-clean-arch"]