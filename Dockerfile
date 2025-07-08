FROM golang:1.24

# Imposta la directory di lavoro
WORKDIR /app


COPY go.mod go.sum ./
RUN go mod download

# Copia i file del progetto
COPY . .

RUN go build -o main .

EXPOSE 8080

CMD ["./main"]