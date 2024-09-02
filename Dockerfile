FROM golang:1.23.0-bullseye AS build-stage

WORKDIR /app

COPY ./ ./

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o /server ./cmd/main.go

FROM scratch AS build-release-stage

WORKDIR /app

COPY --from=build-stage /server /server

EXPOSE 8080

ENTRYPOINT [ "/server" ]