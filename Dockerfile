FROM alpine

COPY bin/pi-sensor-server /app/

WORKDIR /app

ENTRYPOINT ["/app/pi-app-deployer-server"]
