SELECT vul_id as id,
       MAX(vul_severity) as severity,
       COUNT(DISTINCT img_id) as imageCount
FROM kun.cve
GROUP BY vul_id
ORDER BY %s
LIMIT ?, ?
