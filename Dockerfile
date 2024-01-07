# Go build image
FROM golang:1.21.5-bookworm AS go_builder
COPY . streambench
WORKDIR streambench
RUN go install -v ./dupbenchpublisher ./dupbenchsubscriber ./dupbenchtickpublish && \
  go build --race -v -o /go/bin/dupbenchpublisher-race ./dupbenchpublisher && \
  go build --race -v -o /go/bin/dupbenchsubscriber-race ./dupbenchsubscriber


# Subscriber image
FROM gcr.io/distroless/base-nossl-debian12:nonroot AS subscriber-race
COPY --from=go_builder /go/bin/dupbenchsubscriber-race /
ENTRYPOINT ["/dupbenchsubscriber-race"]

FROM gcr.io/distroless/base-nossl-debian12:nonroot AS subscriber
COPY --from=go_builder /go/bin/dupbenchsubscriber /
ENTRYPOINT ["/dupbenchsubscriber"]


# Publisher image
FROM gcr.io/distroless/base-nossl-debian12:nonroot AS publisher-race
COPY --from=go_builder /go/bin/dupbenchpublisher-race /
ENTRYPOINT ["/dupbenchpublisher-race"]

FROM gcr.io/distroless/base-nossl-debian12:nonroot AS publisher
COPY --from=go_builder /go/bin/dupbenchpublisher /
ENTRYPOINT ["/dupbenchpublisher"]

# Tick publisher
FROM gcr.io/distroless/base-nossl-debian12:nonroot AS tickpublish
COPY --from=go_builder /go/bin/dupbenchtickpublish /
ENTRYPOINT ["/dupbenchtickpublish"]
