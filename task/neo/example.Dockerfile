# optional
# --------

FROM ghcr.io/kushsharma/optimus-task-neo-executor:latest

# path to optimus release tar.gz
ARG OPTIMUS_RELEASE_URL

RUN apt install curl tar -y
RUN mkdir -p /opt
RUN curl -sL ${OPTIMUS_RELEASE_URL} | tar xvz optimus
RUN mv optimus /opt/optimus || true
RUN chmod +x /opt/optimus

# or copy like this
COPY task/neo/example.entrypoint.sh /opt/entrypoint.sh
RUN chmod +x /opt/entrypoint.sh

ENTRYPOINT ["/opt/entrypoint.sh"]
CMD ["python3", "/opt/main.py"]