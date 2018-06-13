FROM scratch
WORKDIR /app
ADD build/main /app/
ENTRYPOINT ["./main"]