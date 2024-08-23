FROM golang:latest

WORKDIR /app

COPY *.mod *.sum ./


RUN go install github.com/air-verse/air@latest

COPY . .

EXPOSE 9000

CMD [ "air" ]