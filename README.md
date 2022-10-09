# log

This is an opinionated log package that caters to how I think logging
should happen. The functionalities that are being included are:

- basic configuration
- structured logging capabilities
- annotating PII in log fields and ability to apply masking or other means to it

# Base parameters

- Timestamp format: RFC 3339
- Log format: JSON
- Level key: "lvl"
- Message key: "msg"
- Timestamp key: "ts"
- Available modes for dealing with PII: none, hash, remove
