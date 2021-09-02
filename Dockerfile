FROM golang:1.7.5

WORKDIR /go/src/github.com/netsells/katsu
COPY . .

RUN GIT_COMMIT=$(git rev-list -1 HEAD) && \
  go build -ldflags "-X main.GitCommit=$GIT_COMMIT"

CMD ["./katsu"]
