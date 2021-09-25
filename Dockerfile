FROM golang:1.17-alpine AS build

ENV PORT 10000

COPY . /go/src/ondemand-service
WORKDIR /go/src/ondemand-service

RUN go build -o /go/bin/ondemand-service

FROM alpine
COPY --from=build /go/bin/ondemand-service /go/bin/ondemand-service
CMD [ "/go/bin/ondemand-service" ]