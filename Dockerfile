FROM golang:bullseye

RUN mkdir /app
COPY . /app
WORKDIR /app

RUN go build -o server .

EXPOSE 3000

CMD [ "/app/server" ]


