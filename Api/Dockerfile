FROM golang

WORKDIR /api

COPY . .

RUN go mod tidy

EXPOSE 6069

CMD [ "go", "run", "main.go" ]