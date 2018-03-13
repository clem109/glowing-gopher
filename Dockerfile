FROM iron/base
WORKDIR /app

# copy binary and config into image
COPY config.yml /app/
COPY main /app/
ENTRYPOINT ["./main"]
EXPOSE 8080
