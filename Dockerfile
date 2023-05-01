FROM --platform=amd64 golang:alpine AS build

WORKDIR /build

COPY go.mod whalesayer.go /build/

RUN go build .

RUN ls && pwd

FROM docker/whalesay 

ENV COWPATH "/usr/local/share/cows/"

COPY --from=build /build/whalesay .

CMD "./whalesayer"