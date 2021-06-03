#!/bin/sh
# wait for few seconds to prepare scheduler for the run
sleep 5

# get resources
echo "-- initializing opctl assets"
OPTIMUS_ADMIN=1 /opt/opctl admin build instance $JOB_NAME --project $PROJECT --output-dir $JOB_DIR --type $TASK_TYPE --name $TASK_NAME --scheduled-at $SCHEDULED_AT --host $OPTIMUS_HOSTNAME

# TODO: this doesnt support using back quote sign in env vars, fix it
echo "-- exporting env"
set -o allexport
source $JOB_DIR/in/.env
set +o allexport

echo "-- current envs"
printenv

echo "-- running unit"
exec "$@"