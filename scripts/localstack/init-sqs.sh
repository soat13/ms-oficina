#!/usr/bin/env bash
set -euo pipefail

export AWS_ACCESS_KEY_ID="${AWS_ACCESS_KEY_ID:-000000000000}"
export AWS_SECRET_ACCESS_KEY="${AWS_SECRET_ACCESS_KEY:-000000000000}"
export AWS_DEFAULT_REGION="${AWS_DEFAULT_REGION:-us-east-1}"

ACCOUNT_ID="000000000000"
REGION="${AWS_DEFAULT_REGION}"
MAIN_RETENTION="345600"    # 4 days
DLQ_RETENTION="1209600"    # 14 days
VISIBILITY_TIMEOUT="30"
MAX_RECEIVE_COUNT="5"

QUEUES=(
  "repairorder-diagnostics-finished:false:true"
  "repairorder-canceled:false:true"
  "repairorder-finished:false:true"
  "estimate-created:false:true"
  "estimate-approved:false:true"
  "estimate-rejected:false:true"
  "estimate-canceled:false:true"
  "product-stock-reduce-confirmed:false:true"
  "product-stock-insufficient-detected:false:true"
  "payment-status-changed:true:false"
)

for entry in "${QUEUES[@]}"; do
  IFS=':' read -r base fifo create_dlq <<< "${entry}"

  if [[ "${fifo}" == "true" ]]; then
    queue_name="${base}.fifo"
    dlq_name="${base}-dlq.fifo"
    fifo_attrs=',"FifoQueue":"true","ContentBasedDeduplication":"true"'
  else
    queue_name="${base}"
    dlq_name="${base}-dlq"
    fifo_attrs=""
  fi

  redrive_attrs=""
  if [[ "${create_dlq}" == "true" ]]; then
    dlq_arn="arn:aws:sqs:${REGION}:${ACCOUNT_ID}:${dlq_name}"
    awslocal sqs create-queue \
      --queue-name "${dlq_name}" \
      --attributes "{\"MessageRetentionPeriod\":\"${DLQ_RETENTION}\"${fifo_attrs}}" >/dev/null
    redrive_attrs=",\"RedrivePolicy\":\"{\\\"deadLetterTargetArn\\\":\\\"${dlq_arn}\\\",\\\"maxReceiveCount\\\":\\\"${MAX_RECEIVE_COUNT}\\\"}\""
  fi

  awslocal sqs create-queue \
    --queue-name "${queue_name}" \
    --attributes "{\"MessageRetentionPeriod\":\"${MAIN_RETENTION}\",\"VisibilityTimeout\":\"${VISIBILITY_TIMEOUT}\"${fifo_attrs}${redrive_attrs}}" >/dev/null

  if [[ "${create_dlq}" == "true" ]]; then
    echo "Created queue: ${queue_name} (+ ${dlq_name})"
  else
    echo "Created queue: ${queue_name}"
  fi
done

echo "All SQS queues created."

touch /tmp/sqs-init.done
