# Go build image
FROM golang:1.13.4-buster AS go_builder
COPY . /go/src/streambench
WORKDIR /go/src/streambench
RUN go install --mod=readonly -v ./dupbenchpublisher

# Runtime image
FROM gcr.io/distroless/base-debian10
COPY --from=go_builder /go/bin/dupbenchpublisher /

ENTRYPOINT ["/dupbenchpublisher"]
USER nonroot
