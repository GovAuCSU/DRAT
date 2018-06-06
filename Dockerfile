FROM alpine:latest
COPY bin/drat /
ENV VCAP_SERVICES '{"postgres": [{"credentials": {"username": "dratpg", "host": "dratpg", "password": "dratpg", "name": "dratpg", "port": 5432}, "tags": ["postgres"]}]}'
ENV VCAP_APPLICATION '{}'
ENTRYPOINT ["sleep 10;/drat"]