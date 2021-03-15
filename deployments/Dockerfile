# Initial stage: download modules
FROM golang:1.15 as modules

ARG ssh_key

RUN mkdir /root/.ssh
ADD id_rsa /root/.ssh/id_rsa
ADD go.mod go.sum /m/

###
### *****| sequre |*****
###

# Intermediate stage: Build the binary
FROM golang:1.15 as builder

COPY --from=modules /go/pkg /go/pkg

# add a non-privileged user
RUN useradd -u 10001 person_key_saver

RUN mkdir -p /person_key_saver
ADD . /person_key_saver
WORKDIR /person_key_saver

# Build the binary with go build
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 \
    go build -o ./bin/app ./cmd/person_key_saver/main.go

# Final stage: Run the binary
FROM scratch

# don't forget /etc/passwd from previous stage
COPY --from=builder /etc/passwd /etc/passwd
USER person_key_saver

# and finally the binary
COPY --from=builder /person_key_saver/bin/app /app
COPY --from=builder /person_key_saver/configs/app.toml /configs/app.toml
CMD ["/app"]
