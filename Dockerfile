FROM golang:1.23.6 AS build

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /app ./cmd/worker

FROM ubuntu:noble

RUN apt-get update && apt-get install lsof iputils-ping -y

COPY --from=build /app /app

EXPOSE 8080
CMD ["/app"]
