FROM golang:1.24.2@sha256:30baaea08c5d1e858329c50f29fe381e9b7d7bced11a0f5f1f69a1504cdfbf5e
WORKDIR /app
COPY . .
RUN go build .
EXPOSE 8081
RUN useradd -U -u 1000 appuser && \
    chown -R 1000:1000 /app
USER 1000
CMD [ "./social-network" ]