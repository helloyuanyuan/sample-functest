FROM golang:alpine3.14

ENV LANG=C.UTF-8 CGO_ENABLED=0

RUN apk add --no-cache bash && rm -rf /var/cache/apk/*

WORKDIR /app
COPY . .
ENTRYPOINT [ "/bin/bash" ]

## docker build -t functest:golang .
## docker run -it --name functest --network=sample-functest_functest functest:golang ./buildtest.sh env.prod "$COPIED_TOKEN"