FROM scratch

ENV SYSTEM_LINK_CONFIG=/config/config.yaml
ENV SYSTEM_LINK_LOGS_DIR=/logs

COPY system-link /
ENTRYPOINT ["/system-link"]
