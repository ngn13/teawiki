FROM golang:1.24.0-alpine as builder

RUN apk add sassc make

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

RUN make RELEASE=1

FROM alpine as runner

RUN adduser -h /tw -D -u 1001 runner
RUN chown -R runner:runner /tw

WORKDIR /tw
USER runner

COPY --from=builder /src/locale      ./locale
COPY --from=builder /src/static      ./static
COPY --from=builder /src/views       ./views
COPY --from=builder /src/teawiki.elf ./

ENTRYPOINT ["/tw/teawiki.elf"]
