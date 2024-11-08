FROM golang:1.23 AS build

COPY . /app

WORKDIR /app

RUN go get -v

RUN go build -v -o hf-provisioner-digitalocean

FROM golang:1.23 AS run

COPY --from=build /app/hf-provisioner-digitalocean /hf-provisioner-digitalocean

WORKDIR /

CMD ["/hf-provisioner-digitalocean"]