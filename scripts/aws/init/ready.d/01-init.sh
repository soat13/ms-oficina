#!/usr/bin/env bash
set -euo pipefail

rm -f /tmp/sqs-init.done

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
  "repairorder-diagnostics-finished:fifo=false:dlq=true"
  "repairorder-canceled:fifo=false:dlq=true"
  "estimate-created:fifo=false:dlq=true"
  "estimate-approved:fifo=false:dlq=true"
  "estimate-rejected:fifo=false:dlq=true"
  "estimate-canceled:fifo=false:dlq=true"
  "product-stock-insufficient-detected:fifo=false:dlq=true"
  "estimate-product-stock-reduction-confirmed:fifo=false:dlq=true"
  "repairorder-product-stock-reduction-confirmed:fifo=false:dlq=true"
  "payment-request:fifo=false:dlq=true"
  "payment-status-changed:fifo=true:dlq=false"
)

create_queue_if_not_exists() {
  local queue_name="$1"
  local attributes="$2"

  if awslocal sqs get-queue-url \
    --queue-name "${queue_name}" \
    --region "${REGION}" >/dev/null 2>&1; then
    echo "Queue ${queue_name} already exists, skipping..."
  else
    awslocal sqs create-queue \
      --queue-name "${queue_name}" \
      --attributes "${attributes}" \
      --region "${REGION}" >/dev/null
    echo "Queue ${queue_name} created"
  fi
}

echo "Checking/creating SQS queues..."

for entry in "${QUEUES[@]}"; do
  IFS=':' read -r base fifo_kv dlq_kv <<< "${entry}"
  fifo="${fifo_kv#fifo=}"
  create_dlq="${dlq_kv#dlq=}"

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
    create_queue_if_not_exists \
      "${dlq_name}" \
      "{\"MessageRetentionPeriod\":\"${DLQ_RETENTION}\"${fifo_attrs}}"
    redrive_attrs=",\"RedrivePolicy\":\"{\\\"deadLetterTargetArn\\\":\\\"${dlq_arn}\\\",\\\"maxReceiveCount\\\":\\\"${MAX_RECEIVE_COUNT}\\\"}\""
  fi

  create_queue_if_not_exists \
    "${queue_name}" \
    "{\"MessageRetentionPeriod\":\"${MAIN_RETENTION}\",\"VisibilityTimeout\":\"${VISIBILITY_TIMEOUT}\"${fifo_attrs}${redrive_attrs}}"
done

STOCK_REDUCTION_CONFIRMED_TOPIC_NAME="product-stock-reduction-confirmed"
STOCK_REDUCTION_CONFIRMED_SUBSCRIBERS=(
  "estimate-product-stock-reduction-confirmed"
  "repairorder-product-stock-reduction-confirmed"
)

echo "Checking/creating SNS topic..."

STOCK_REDUCTION_CONFIRMED_TOPIC_ARN=$(awslocal sns create-topic \
  --name "${STOCK_REDUCTION_CONFIRMED_TOPIC_NAME}" \
  --region "${REGION}" \
  --query 'TopicArn' \
  --output text)

echo "Topic ready: ${STOCK_REDUCTION_CONFIRMED_TOPIC_ARN}"

subscribe_queue_to_topic() {
  local queue_name="$1"
  local topic_arn="$2"

  local queue_url
  queue_url=$(awslocal sqs get-queue-url \
    --queue-name "${queue_name}" \
    --region "${REGION}" \
    --query 'QueueUrl' \
    --output text)

  local queue_arn
  queue_arn=$(awslocal sqs get-queue-attributes \
    --queue-url "${queue_url}" \
    --attribute-names QueueArn \
    --region "${REGION}" \
    --query 'Attributes.QueueArn' \
    --output text)

  local policy
  policy=$(cat <<EOF
{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Principal":"*","Action":"sqs:SendMessage","Resource":"${queue_arn}","Condition":{"ArnEquals":{"aws:SourceArn":"${topic_arn}"}}}]}
EOF
)

  local escaped_policy
  escaped_policy=$(printf '%s' "$policy" | python3 -c 'import json,sys; print(json.dumps(sys.stdin.read()))')

  awslocal sqs set-queue-attributes \
    --queue-url "${queue_url}" \
    --attributes "{\"Policy\":${escaped_policy}}" \
    --region "${REGION}"

  local existing_subscription_arn
  existing_subscription_arn=$(awslocal sns list-subscriptions-by-topic \
    --topic-arn "${topic_arn}" \
    --region "${REGION}" \
    --query "Subscriptions[?Endpoint=='${queue_arn}'].SubscriptionArn | [0]" \
    --output text)

  if [ "${existing_subscription_arn}" = "None" ] || [ -z "${existing_subscription_arn}" ]; then
    awslocal sns subscribe \
      --topic-arn "${topic_arn}" \
      --protocol sqs \
      --notification-endpoint "${queue_arn}" \
      --region "${REGION}" >/dev/null
    echo "Queue ${queue_name} subscribed to topic ${STOCK_REDUCTION_CONFIRMED_TOPIC_NAME}"
  else
    echo "Subscription for ${queue_name} already exists, skipping..."
  fi
}

for subscriber in "${STOCK_REDUCTION_CONFIRMED_SUBSCRIBERS[@]}"; do
  subscribe_queue_to_topic "${subscriber}" "${STOCK_REDUCTION_CONFIRMED_TOPIC_ARN}"
done

echo "Done."
echo "Resources ready:"
for entry in "${QUEUES[@]}"; do
  IFS=':' read -r base fifo_kv dlq_kv <<< "${entry}"
  fifo="${fifo_kv#fifo=}"
  if [[ "${fifo}" == "true" ]]; then
    echo "- SQS queue: ${base}.fifo"
  else
    echo "- SQS queue: ${base}"
  fi
done
echo "- SNS topic: ${STOCK_REDUCTION_CONFIRMED_TOPIC_NAME}"

touch /tmp/sqs-init.done
