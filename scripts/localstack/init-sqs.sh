#!/bin/bash
set -euo pipefail

QUEUES=(
  "repairorder-diagnostics-finished"
  "repairorder-canceled"
  "repairorder-finished"
  "estimate-created"
  "estimate-approved"
  "estimate-rejected"
  "product-stock-reduce-confirmed"
  "product-stock-insufficient-detected"
  "payment-status-changed"
)

for q in "${QUEUES[@]}"; do
  awslocal sqs create-queue --queue-name "${q}-dlq"

  DLQ_ARN=$(awslocal sqs get-queue-attributes \
    --queue-url "http://localhost:4566/000000000000/${q}-dlq" \
    --attribute-names QueueArn \
    --query 'Attributes.QueueArn' --output text)

  awslocal sqs create-queue --queue-name "$q" \
    --attributes "{\"RedrivePolicy\":\"{\\\"deadLetterTargetArn\\\":\\\"${DLQ_ARN}\\\",\\\"maxReceiveCount\\\":\\\"3\\\"}\"}"

  echo "Created queue: $q (+ DLQ)"
done

echo "All SQS queues created."
