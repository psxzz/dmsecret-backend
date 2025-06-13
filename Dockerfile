FROM golang:1.23

COPY --chmod=755 bin/ghosty_link .

CMD ["./ghosty_link"]
