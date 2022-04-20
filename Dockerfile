FROM golang:alpine3.14

ENV LANG=C.UTF-8

RUN apk add --no-cache bash && rm -rf /var/cache/apk/*

WORKDIR /app
COPY . .
ENTRYPOINT [ "/bin/bash" ]

## docker build -t functest:golang .
## docker run -it --name functest functest:golang ./buildtest.sh prod