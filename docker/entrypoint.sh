#!/usr/bin/env sh
exec /kun-api \
  --mongodb-addr "${MONGODB_ADDR}" \
  --mongodb-user "${MONGODB_USER}" \
  --mongodb-pass "${MONGODB_PASS}"
