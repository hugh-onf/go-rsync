FROM golang AS builder
WORKDIR /src
RUN git clone https://github.com/hugh-onf/go-rsync.git && cd /src/go-rsync && go build

FROM ubuntu AS runner
COPY --from=builder /src/go-rsync/go-rsync /usr/bin
RUN chmod +x /usr/bin/go-rsync
ENTRYPOINT [ "go-rsync" ]