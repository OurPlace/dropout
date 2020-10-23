# For those that don't know, because it seems common enough that people don't,
# you can have multi staged dockerfiles that specify a builder and the running
# image.
#
# Read more here:
# https://docs.docker.com/develop/develop-images/multistage-build/
#
FROM golang:1.14-alpine AS builder

# Update our certificates incase our builder/application image is out of date (it is).
RUN apk --no-cache add ca-certificates upx git build-base gcc abuild binutils binutils-doc gcc-doc

# Make output directory.
RUN mkdir /.bin

WORKDIR /go/src/github.com/ourplace/dropout

COPY . .

RUN GOOS=linux GARCH=amd64 \
    go build \
    -o /.bin/dropout \
    -ldflags  "-s -w -installsuffix nocgo -linkmode external -extldflags -static" \
    ./cmd/dropout/main.go

# If we want to make dropout smaller, we can always throw it into upx c:
# Stub in UPX here I guess
# RUN upx --lzma /.bin/dropout


# Main Image
# TODO:
# - Figure out if we have a policy against scratch images, or if we want to use
# something specific, or configured in a certain way.
FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /.bin/dropout /dropout

CMD ["/dropout"]