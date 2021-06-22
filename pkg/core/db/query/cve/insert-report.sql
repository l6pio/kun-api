INSERT INTO kun.cve (img_id, img_name, img_size,
                     art_name, art_version, art_type, art_lang, art_licenses, art_cpes, art_purl,
                     vul_id, vul_data_source, vul_namespace, vul_severity, vul_urls, vul_desc,
                     vul_cvss_version, vul_cvss_base_score, vul_cvss_exploit_score, vul_cvss_impact_score,
                     vul_fix_versions, vul_fix_state, timestamp)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
