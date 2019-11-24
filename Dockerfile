# Go build image
FROM golang:1.13.4-buster AS go_builder
COPY . /go/src/streambench
WORKDIR /go/src/streambench
RUN go install --mod=readonly -v ./dupbenchpublisher ./dupbenchsubscriber

# Subscriber image
FROM gcr.io/distroless/base-debian10 AS subscriber
COPY --from=go_builder /go/bin/dupbenchsubscriber /
ENTRYPOINT ["/dupbenchsubscriber"]
USER nonroot

# Publisher image
FROM gcr.io/distroless/base-debian10 AS publisher
COPY --from=go_builder /go/bin/dupbenchpublisher /
ENTRYPOINT ["/dupbenchpublisher"]
USER nonroot
