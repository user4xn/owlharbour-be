FROM golang:1.19-bullseye
ENV TZ=Asia/Jakarta

COPY . /app/
WORKDIR /app/

RUN go build -v .

CMD ["./owlharbour-api"]