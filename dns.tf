resource "aws_api_gateway_domain_name" "emojibot" {
  domain_name              = local.current.dns_name
  regional_certificate_arn = aws_acm_certificate_validation.emojiBot.certificate_arn

  endpoint_configuration {
    types = ["REGIONAL"]
  }
}


// DNS record for emojiBot
resource "aws_route53_record" "emojibot" {
  name    = aws_api_gateway_domain_name.emojibot.domain_name
  type    = "A"
  zone_id = data.aws_route53_zone.ctrlokDev.id

  alias {
    evaluate_target_health = true
    name                   = aws_api_gateway_domain_name.emojibot.regional_domain_name
    zone_id                = aws_api_gateway_domain_name.emojibot.regional_zone_id
  }
}

resource "aws_acm_certificate" "emojiBot" {
  domain_name       = "emojibot.aws.ctrlok.dev"
  validation_method = "DNS"
}

data "aws_route53_zone" "ctrlokDev" {
  name         = "aws.ctrlok.dev."
  private_zone = false
}

resource "aws_route53_record" "emojiBot_cert_validation" {
  name    = aws_acm_certificate.emojiBot.domain_validation_options[0].resource_record_name
  type    = aws_acm_certificate.emojiBot.domain_validation_options[0].resource_record_type
  zone_id = data.aws_route53_zone.ctrlokDev.id
  records = [
  aws_acm_certificate.emojiBot.domain_validation_options[0].resource_record_value]
  ttl = 60
}

resource "aws_acm_certificate_validation" "emojiBot" {
  certificate_arn = aws_acm_certificate.emojiBot.arn
  validation_record_fqdns = [
  aws_route53_record.emojiBot_cert_validation.fqdn]
}
