resource "aws_acm_certificate" "api_gateway_domain" {
  domain_name       = "api.${var.aws_route53_root_domain}"
  validation_method = "DNS"

  validation_option {
    domain_name       = "api.${var.aws_route53_root_domain}"
    validation_domain = var.aws_route53_root_domain
  }
}

resource "aws_acm_certificate_validation" "api_gateway_domain" {
  certificate_arn = aws_acm_certificate.api_gateway_domain.arn
}

resource "aws_apigatewayv2_domain_name" "main" {
  depends_on = [aws_acm_certificate_validation.api_gateway_domain]

  domain_name = "api.${var.aws_route53_root_domain}"

  domain_name_configuration {
    certificate_arn = aws_acm_certificate.api_gateway_domain.arn
    endpoint_type   = "REGIONAL"
    security_policy = "TLS_1_2"
  }
}

resource "aws_apigatewayv2_api_mapping" "prod" {
  api_id = aws_apigatewayv2_api.main.id

  domain_name = aws_apigatewayv2_domain_name.main.id
  stage       = aws_apigatewayv2_stage.prod.id
}

resource "aws_route53_record" "apigateway_validation_record" {
  for_each = {
    for dvo in aws_acm_certificate.api_gateway_domain.domain_validation_options : dvo.domain_name => {
      name   = dvo.resource_record_name
      record = dvo.resource_record_value
      type   = dvo.resource_record_type
    }
  }

  zone_id = var.aws_route53_zoneid

  name    = each.value.name
  type    = each.value.type
  ttl     = 300
  records = [each.value.record]

  allow_overwrite = true
}

resource "aws_route53_record" "apigateway_record" {
  zone_id = var.aws_route53_zoneid

  name    = "api"
  type    = "CNAME"
  ttl     = 300
  records = [aws_apigatewayv2_domain_name.main.domain_name_configuration[0].target_domain_name]

  allow_overwrite = true
}
