FROM golang:1.18-alpine
WORKDIR /app
COPY . .
COPY ./go.* .
RUN go mod download
COPY . .

EXPOSE 8001

ENTRYPOINT [ "go","run", "cmd/driver/main.go" ]