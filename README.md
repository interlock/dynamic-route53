# Dynamic Route 53 DNS

This cli application enables updating an AWS Route53 HostedZone Record with an ipv4/ipv6 address.

## Usage

## Configuration

- AWS_ACCESS_KEY - required env variable to interact with AWS
- AWS_SECRET_KEY - required env variable to interact wiht AWS

```
      --comment string          comment to RecordSet
  -d, --domain string           domain to update
  -f, --frequency int           frequency to refresh. 0 value is once
  -h, --help                    help for dynamic-route53
  -z, --hosted-zone-id string   hosted zone ID
      --immediate               immediately start refresh on start (default true)
  -l, --lookup string           which lookup agent to lookup (default "ifconfig.co")
  -n, --network string          tcp4 or tcp6, sets A or AAAA as appropriate (default "tcp4")
      --profile                 Enable profiling
      --profile-port int32      Port for profiling (default 6660)
      --ttl int                 TTL to set on the Route53 Record (default 3600)
```