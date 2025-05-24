# builds the binary
FROM golang:1.24.0-alpine as build

ARG TARGETOS
ARG TARGETARCH

RUN apk add git sassc make

WORKDIR /src

COPY *.mod    ./
COPY *.sum    ./

RUN go mod download

COPY *.go     ./
COPY Makefile ./
COPY .git     ./.git
COPY config   ./config
COPY consts   ./consts
COPY locale   ./locale
COPY log      ./log
COPY repo     ./repo
COPY routes   ./routes
COPY static   ./static
COPY util     ./util
COPY views    ./views

ENV CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH
RUN make RELEASE=1

# runs the binary
FROM alpine

RUN adduser -h /tw -D -u 1001 runner && \
    chown -R runner:runner /tw

WORKDIR /tw
USER runner

COPY --from=build /src/locale      ./locale
COPY --from=build /src/static      ./static
COPY --from=build /src/views       ./views
COPY --from=build /src/teawiki.elf ./

ENTRYPOINT ["/tw/teawiki.elf"]
