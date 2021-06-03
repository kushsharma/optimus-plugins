# optional
# --------

FROM cosmopod/optimus-task-neo-executor:latest

# path to optimus release tar.gz
ARG OPTIMUS_RELEASE_URL

RUN apt install curl tar -y
RUN mkdir -p /opt
RUN curl -sL ${OPTIMUS_RELEASE_URL} | tar xvz opctl
RUN mv opctl /opt/opctl || true
RUN chmod +x /opt/opctl

# either pack the shell file in a single string
# RUN echo '#!/bin/sh' >> /opt/entrypoint.sh && echo '# wait for few seconds to prepare scheduler for the run' >> /opt/entrypoint.sh && echo 'sleep 5' >> /opt/entrypoint.sh && echo '# get resources' >> /opt/entrypoint.sh && echo 'echo "-- initializing opctl assets"' >> /opt/entrypoint.sh && echo 'OPTIMUS_ADMIN=1 /opt/opctl admin build instance $JOB_NAME --project $PROJECT --output-dir $JOB_DIR --type $TASK_TYPE --name $TASK_NAME --scheduled-at $SCHEDULED_AT --host $OPTIMUS_HOSTNAME' >> /opt/entrypoint.sh && echo '# TODO: this doesnt support using back quote sign in env vars, fix it' >> /opt/entrypoint.sh && echo 'echo "-- exporting env"' >> /opt/entrypoint.sh && echo 'set -o allexport' >> /opt/entrypoint.sh && echo 'source $JOB_DIR/in/.env' >> /opt/entrypoint.sh && echo 'set +o allexport' >> /opt/entrypoint.sh && echo 'echo "-- current envs"' >> /opt/entrypoint.sh && echo 'printenv' >> /opt/entrypoint.sh && echo 'echo "-- running unit"' >> /opt/entrypoint.sh && echo 'exec "$@"' >> /opt/entrypoint.sh

# or copy like this
COPY task/neo/example.entrypoint.sh /opt/entrypoint.sh
RUN chmod +x /opt/entrypoint.sh

ENTRYPOINT ["/opt/entrypoint.sh"]
CMD ["python3", "/opt/main.py"]