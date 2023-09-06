resource "random_pet" "website_bucket" {
  prefix = "games4l-website"
  length = 2
}

resource "aws_s3_bucket" "website_bucket" {
  bucket = random_pet.website_bucket.id
}

resource "aws_s3_bucket_public_access_block" "website_bucket" {
  bucket = aws_s3_bucket.website_bucket.id

  block_public_acls       = false
  block_public_policy     = false
  ignore_public_acls      = false
  restrict_public_buckets = false
}

resource "aws_s3_bucket_website_configuration" "website_bucket" {
  bucket = aws_s3_bucket.website_bucket.id

  index_document {
    suffix = "index.html"
  }

  error_document {
    key = "error.html"
  }
}

# Read access
# Needs to await the bucket creation

resource "time_sleep" "wait_2_seconds" {
  depends_on      = [aws_s3_bucket.website_bucket]
  create_duration = "2s"
}

data "aws_iam_policy_document" "allow_public_access" {
  statement {
    principals {
      type        = "*"
      identifiers = ["*"]
    }

    actions = [
      "s3:GetObject"
    ]

    resources = [
      "${aws_s3_bucket.website_bucket.arn}/*",
    ]
  }
}

resource "aws_s3_bucket_policy" "website_policy" {
  depends_on = [time_sleep.wait_2_seconds]
  bucket     = aws_s3_bucket.website_bucket.id

  policy = data.aws_iam_policy_document.allow_public_access.json
}

resource "aws_s3_object" "files" {
  depends_on = [aws_s3_bucket_policy.website_policy]
  for_each   = fileset("../${path.module}/apps/memories/dist/", "**")
  bucket     = aws_s3_bucket.website_bucket.id

  key          = each.value
  content_type = lookup(local.mime_types, regex("[^.]+$", each.value), null)
  source       = "../${path.module}/apps/memories/dist/${each.value}"
  etag         = filemd5("../${path.module}/apps/memories/dist/${each.value}")
}
