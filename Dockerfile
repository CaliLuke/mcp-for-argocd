FROM golang:1.26 AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/argocd-mcp ./cmd/argocd-mcp

FROM gcr.io/distroless/base-debian12
COPY --from=build /out/argocd-mcp /argocd-mcp
EXPOSE 3000
ENTRYPOINT ["/argocd-mcp"]
CMD ["http", "--port", "3000"]
