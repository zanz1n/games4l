resource "aws_cloudfront_distribution" "website" {
  origin {
    custom_origin_config {
      http_port              = 80
      https_port             = 443
      origin_protocol_policy = "http-only"

      origin_ssl_protocols = ["TLSv1", "TLSv1.1", "TLSv1.2"]
    }

    domain_name = aws_s3_bucket_website_configuration.website_bucket.website_endpoint

    origin_id = var.website_cloudflare_domain
  }

  enabled = true
  # default_root_object = "index.html"

  default_cache_behavior {
    viewer_protocol_policy = "redirect-to-https"
    compress               = true
    allowed_methods        = ["GET", "HEAD"]
    cached_methods         = ["GET", "HEAD"]
    target_origin_id       = var.website_cloudflare_domain
    min_ttl                = 0
    default_ttl            = 1800
    max_ttl                = 86400

    forwarded_values {
      query_string = false
      cookies {
        forward = "none"
      }
    }
  }

  aliases = [var.website_cloudflare_domain]

  restrictions {
    geo_restriction {
      restriction_type = "none"
    }
  }

  viewer_certificate {
    acm_certificate_arn = aws_acm_certificate.website.arn
    ssl_support_method  = "sni-only"
  }
}

resource "aws_acm_certificate" "website" {
  provider          = aws.virginia
  domain_name       = var.website_cloudflare_domain
  validation_method = "DNS"

  validation_option {
    domain_name       = var.website_cloudflare_domain
    validation_domain = var.cloudflare_root_domain
  }
}

resource "aws_acm_certificate_validation" "website" {
  provider        = aws.virginia
  certificate_arn = aws_acm_certificate.website.arn
}

resource "cloudflare_record" "website_validation" {
  for_each = {
    for dvo in aws_acm_certificate.website.domain_validation_options : dvo.domain_name => {
      name   = dvo.resource_record_name
      record = dvo.resource_record_value
      type   = dvo.resource_record_type
    }
  }
  type  = each.value.type
  name  = each.value.name
  value = each.value.record

  zone_id         = var.cloudflare_zoneid
  proxied         = false
  allow_overwrite = true
}

resource "cloudflare_record" "website_record" {
  zone_id = var.cloudflare_zoneid

  type            = "CNAME"
  proxied         = false
  allow_overwrite = true
  name            = var.website_cloudflare_domain
  value           = aws_cloudfront_distribution.website.domain_name
}
