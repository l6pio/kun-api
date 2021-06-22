CREATE TABLE IF NOT EXISTS kun.cve
(
    img_id                 String,
    img_name               String,
    img_size               Int64,
    art_name               String,
    art_version            String,
    art_type               String,
    art_lang               String,
    art_licenses           Array(String),
    art_cpes               Array(String),
    art_purl               String,
    vul_id                 String,
    vul_data_source        String,
    vul_namespace          String,
    vul_severity           Int64,
    vul_urls               Array(String),
    vul_desc               String,
    vul_cvss_version       Float64,
    vul_cvss_base_score    Float64,
    vul_cvss_exploit_score Float64,
    vul_cvss_impact_score  Float64,
    vul_fix_versions       Array(String),
    vul_fix_state          Int64,
    timestamp              DateTime
)
    ENGINE = MergeTree()
        PARTITION BY toYYYYMM(timestamp)
        ORDER BY img_id
