FROM golang:1.22


RUN go version


ENV GOPATH=/

COPY ./ ./


RUN go mod download

RUN go build -o avito-rent ./cmd/main.go
RUN go build -o migrate ./cmd/migrate/main.go


COPY entrypoint.sh ./entrypoint.sh
COPY wait-for-db.sh ./wait-for-db.sh

RUN apt-get update && \
    apt-get install -y postgresql-client 

RUN chmod +x ./entrypoint.sh ./wait-for-db.sh 


ENTRYPOINT ["sh", "/entrypoint.sh"]


CMD ["./avito-rent"]
