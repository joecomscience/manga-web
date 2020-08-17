FROM golang:1.13-alpine as build
ENV GOOS=linux
ENV GOARCH=amd64

WORKDIR /app
COPY . .

RUN go mod download && \
    go build main.go

FROM alpine:3.11 as release

LABEL maintrainer="joewalker"
COPY --from=build /app/main /app
EXPOSE 3000
CMD ["./app"]