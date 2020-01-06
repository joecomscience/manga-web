FROM alpine:3.11

LABEL maintrainer="joewalker"
COPY ./main /app
EXPOSE 3000
CMD ["./app"]