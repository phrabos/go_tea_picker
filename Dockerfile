# syntax=docker/dockerfile:1

FROM golang:1.19
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY *.go ./
COPY *.json ./
RUN CGO_ENABLED=0 GOOS=linux go build -o /tea-picker
CMD ["/tea-picker"]