CREATE TABLE IF NOT EXISTS kun.img
(
    id         String,
    img_id     String,
    img_name   String,
    img_status Int64,
    timestamp  DateTime
)
    ENGINE = MergeTree()
        PARTITION BY toYYYYMM(timestamp)
        ORDER BY img_id
