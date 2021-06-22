SELECT id
FROM kun.img
WHERE img_id = ? AND timestamp > ?
ORDER BY timestamp
