### Node Build (for frontend)
FROM node:latest AS node_builder

WORKDIR /frontend

COPY ./frontend/package*.json ./
RUN npm install

COPY ./frontend .
RUN npm run build

### Go Build (for main app)
FROM golang:1.13 AS go_builder

WORKDIR /go/src/github.com/tjhorner/dash
RUN go get -u github.com/gobuffalo/packr/v2/packr2

COPY . .
RUN go get -d -v .

COPY --from=node_builder /frontend/build ./frontend/build/

RUN packr2 build -v -ldflags="-linkmode external -extldflags -static -s -w" -o /dash *.go

### Packaged (single binary!)
# We use the distroless/static image since it includes a list of CAs and tzinfo, but is also very slim
FROM gcr.io/distroless/static:8bef63d2c8654ff89358430c7df5778162ab6027

EXPOSE 3000

COPY --from=go_builder /dash /dash
CMD ["/dash"]