SELECT img_id                   AS id,
       img_name                 AS name,
       MAX(img_size)            AS size,
       COUNT(DISTINCT art_name) AS artCount,
       COUNT(DISTINCT vul_id)   AS vulCount
FROM kun.cve
GROUP BY id, name
ORDER BY %s
LIMIT ?, ?
