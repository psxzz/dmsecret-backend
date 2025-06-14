FROM golang:1.24

COPY --chmod=755 bin/ghosty_link .

CMD ["./ghosty_link"]
