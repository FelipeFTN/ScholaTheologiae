FROM golang:latest

WORKDIR /app/api
RUN mkdir -p /app/books/
COPY books/* /app/books/
COPY ./api /app/api
RUN make build

CMD ["./bin/schola-theologiae-api"]
