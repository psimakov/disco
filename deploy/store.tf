# dataset and tables for disco

resource "google_bigquery_dataset" "disco" {
  dataset_id    = "disco"
  friendly_name = "disco"
  description   = "disco dataset"
  location      = "US"

  labels = {
    env = "demo"
  }
}

// Licenses

resource "google_bigquery_table" "licenses" {
  dataset_id = google_bigquery_dataset.disco.dataset_id
  table_id   = "licenses"
  deletion_protection = false

  time_partitioning {
    type = "DAY"
  }

  labels = {
    env = "demo"
  }

  schema = data.template_file.schema_licenses.rendered

}


resource "google_bigquery_table" "licenses_day" {
  dataset_id = google_bigquery_dataset.disco.dataset_id
  table_id   = "licenses_day"
  deletion_protection = false

  view {
    use_legacy_sql = false
    query          = <<EOF
SELECT  	
  image,
  sha,
  UPPER(CASE WHEN CONTAINS_SUBSTR(name, '-') THEN SPLIT(name, '-')[SAFE_OFFSET(0)] ELSE SUBSTR(name, 0, 3) END) short_name,
  name,
  package,
  FORMAT_TIMESTAMP("%Y-%m-%d",updated) updated
FROM `${google_bigquery_table.licenses.project}.${google_bigquery_table.licenses.dataset_id}.${google_bigquery_table.licenses.table_id}` 
WHERE _PARTITIONTIME IS NULL OR _PARTITIONTIME > TIMESTAMP_TRUNC(TIMESTAMP_SUB(CURRENT_TIMESTAMP, INTERVAL 1 DAY), DAY)
GROUP BY 1,2,3,4,5,6
EOF
  }
}

// Vulnerabilities

resource "google_bigquery_table" "vulnerabilities" {
  dataset_id = google_bigquery_dataset.disco.dataset_id
  table_id   = "vulnerabilities"
  deletion_protection = false

  time_partitioning {
    type = "DAY"
  }

  labels = {
    env = "demo"
  }

  schema = data.template_file.schema_vulnerabilities.rendered
}


resource "google_bigquery_table" "vulnerabilities_day" {
  dataset_id = google_bigquery_dataset.disco.dataset_id
  table_id   = "vulnerabilities_day"
  deletion_protection = false

  view {
    use_legacy_sql = false
    query          = <<EOF
SELECT  	
  image,
  sha,
  cve,
  severity,
  package,
  version,
  title,
  description,
  url,
  FORMAT_TIMESTAMP("%Y-%m-%d",updated) updated
FROM `${google_bigquery_table.vulnerabilities.project}.${google_bigquery_table.vulnerabilities.dataset_id}.${google_bigquery_table.vulnerabilities.table_id}` 
WHERE _PARTITIONTIME IS NULL OR _PARTITIONTIME > TIMESTAMP_TRUNC(TIMESTAMP_SUB(CURRENT_TIMESTAMP, INTERVAL 1 DAY), DAY)
GROUP BY 1,2,3,4,5,6,7,8,9,10 
EOF
  }
}

// Packages

resource "google_bigquery_table" "packages" {
  dataset_id = google_bigquery_dataset.disco.dataset_id
  table_id   = "packages"
  deletion_protection = false

  time_partitioning {
    type = "DAY"
  }

  labels = {
    env = "demo"
  }

  schema = data.template_file.schema_packages.rendered
}

resource "google_bigquery_table" "packages_day" {
  dataset_id = google_bigquery_dataset.disco.dataset_id
  table_id   = "packages_day"
  deletion_protection = false

  view {
    use_legacy_sql = false
    query          = <<EOF
SELECT  	
  image,
  sha,
  format,
  provider,
  originator,
  package,
  version,
  source,
  license,
  FORMAT_TIMESTAMP("%Y-%m-%d",updated) updated
FROM `${google_bigquery_table.packages.project}.${google_bigquery_table.packages.dataset_id}.${google_bigquery_table.packages.table_id}` 
WHERE _PARTITIONTIME IS NULL OR _PARTITIONTIME > TIMESTAMP_TRUNC(TIMESTAMP_SUB(CURRENT_TIMESTAMP, INTERVAL 1 DAY), DAY)
GROUP BY 1,2,3,4,5,6,7,8,9,10 
EOF
  }
}
