terraform {
  backend "gcs" {
    prefix = "rusuden"
  }
}

variable "project" {}

variable "MG_SIGNING_KEY" {}
variable "MG_API_KEY" {}
variable "MG_SENDER" {}
variable "MG_RECIPIENT" {}

locals {
  name   = "rusuden"
  region = "us-central1"
}

provider "google" {
  project = var.project
  region  = local.region
}

provider "mailgun" {
  api_key = var.MG_API_KEY
}

data "archive_file" "archive" {
  type        = "zip"
  output_path = "source.zip"

  source_dir = "."
  excludes   = split("\n", file(".gitignore"))
}

resource "google_storage_bucket" "bucket" {
  name     = "${var.project}-${local.name}"
  location = local.region
}

resource "google_storage_bucket_object" "object" {
  name   = "${data.archive_file.archive.output_md5}.zip"
  bucket = google_storage_bucket.bucket.name
  source = data.archive_file.archive.output_path
}

resource "google_cloudfunctions_function" "function" {
  name = local.name

  runtime             = "go113"
  available_memory_mb = 512
  entry_point         = "Handle"
  trigger_http        = true

  source_archive_bucket = google_storage_bucket.bucket.name
  source_archive_object = google_storage_bucket_object.object.name

  environment_variables = {
    MG_SIGNING_KEY = var.MG_SIGNING_KEY,
    MG_API_KEY     = var.MG_API_KEY,
    MG_SENDER      = var.MG_SENDER,
    MG_RECIPIENT   = var.MG_RECIPIENT,
  }
}

resource "mailgun_route" "route" {
  priority   = "10"
  expression = "match_header(\"Subject\", \"【SMARTalk】メッセージお預かり通知\")"
  actions    = ["forward(\"${google_cloudfunctions_function.function.https_trigger_url}\")"]
}
