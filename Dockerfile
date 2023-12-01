FROM alpine:3.16-latest


RUN go build .


ENTRYPOINT [ "bash","./webscrapping-go.exe" ]