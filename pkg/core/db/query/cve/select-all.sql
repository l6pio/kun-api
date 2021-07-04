SELECT vul_id                 AS id,
       MAX(vul_severity)      AS severity,
       MAX(vul_fix_state)     AS fixState,
       COUNT(DISTINCT img_id) AS imageCount
FROM kun.cve
GROUP BY id
ORDER BY %s
LIMIT ?, ?
