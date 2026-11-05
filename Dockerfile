FROM golang:1.27

COPY api/public/api.yaml .
COPY --chmod=755 bin/ghosty_link .

CMD ["./ghosty_link"]
