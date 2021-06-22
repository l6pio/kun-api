SELECT vul_id                 AS id,
       MAX(vul_severity)      AS severity,
       COUNT(DISTINCT img_id) AS imageCount
FROM kun.cve
GROUP BY vul_id
ORDER BY %s
LIMIT ?, ?
