---
title: "Debugging TLS/SSL Issues with OpenSSL"
slug: a3f7c219
aliases: openssl-tls-debugging
date: 2026-03-20
tags: [security, networking, tools]
lang: en
draft: false
---

TLS failures are common and often cryptic: a curl error, a rejected handshake, or a certificate that expires silently at 3 AM. OpenSSL ships with several subcommands that make most of these problems diagnosable without external tools.

## Fun Facts

**Fact 1.** The `s_client` subcommand was originally intended as a diagnostic tool, not a production client. It has been part of OpenSSL since version 0.9.5 (released in 2000), making it older than most Linux distributions in common use today.

**Fact 2.** PEM format is Base64-encoded DER with a header line. The name "PEM" stands for Privacy Enhanced Mail, a 1993 email encryption standard that was never widely adopted but whose encoding survived in TLS infrastructure.

**Fact 3.** TLS alert codes are standardized in RFC 8446. Alert code 42 means `bad_certificate`; alert 48 means `unknown_ca`. These numeric codes appear in raw packet captures and in some OpenSSL error messages.

---

## Tips and Tricks

### 1. Inspect the Full Certificate Chain with s_client

`s_client` connects to a TLS endpoint and prints the negotiated chain, the certificate details, and the handshake result. The `-showcerts` flag prints every certificate the server sends, not just the leaf.

```bash
openssl s_client -connect example.com:443 -showcerts </dev/null 2>/dev/null
```

To test a specific SNI name on a shared-IP host:

```bash
openssl s_client -connect 93.184.216.34:443 \
  -servername example.com \
  -showcerts </dev/null 2>/dev/null
```

The output includes `Verify return code: 0 (ok)` when the chain validates against the system trust store. Any nonzero code is an error worth investigating.

### 2. Read Certificate Details with x509

Once you have a PEM certificate (either from a file or piped from `s_client`), `openssl x509` decodes its fields.

```bash
# Print subject, issuer, validity dates, and SANs
openssl x509 -in cert.pem -noout -text | grep -E \
  "Subject:|Issuer:|Not Before|Not After|DNS:"
```

To read a DER-encoded certificate (common in Java keystores and Windows exports):

```bash
openssl x509 -in cert.der -inform DER -noout -text
```

To extract just the expiry date as a machine-readable string:

```bash
openssl x509 -in cert.pem -noout -enddate
# notAfter=Mar 20 12:00:00 2027 GMT
```

### 3. Test Specific Cipher Suites

Some servers reject clients that offer deprecated ciphers; others refuse to negotiate modern ones. You can force a specific cipher for the handshake with `-cipher` (TLS 1.2 and below) or `-ciphersuites` (TLS 1.3).

```bash
# Test with a specific TLS 1.2 cipher
openssl s_client -connect example.com:443 \
  -tls1_2 \
  -cipher ECDHE-RSA-AES256-GCM-SHA384 \
  </dev/null 2>&1 | grep -E "Cipher|Alert|error"

# Test TLS 1.3 cipher suite
openssl s_client -connect example.com:443 \
  -tls1_3 \
  -ciphersuites TLS_AES_256_GCM_SHA384 \
  </dev/null 2>&1 | grep -E "Cipher|Alert|error"
```

If the handshake fails with `no ciphers available` or `handshake failure`, the server does not support that combination.

### 4. Monitor Certificate Expiry with a Bash Script

A simple script that checks one or more hosts and exits nonzero when any certificate expires within a threshold. Wire this into a cron job or a monitoring system.

```bash
#!/usr/bin/env bash
# check-cert-expiry.sh
# Usage: ./check-cert-expiry.sh example.com:443 api.example.com:443
# Exit 1 if any cert expires within WARN_DAYS days.

WARN_DAYS="${WARN_DAYS:-30}"
STATUS=0

for HOST_PORT in "$@"; do
  HOST="${HOST_PORT%%:*}"
  PORT="${HOST_PORT##*:}"

  EXPIRY=$(openssl s_client \
    -connect "${HOST_PORT}" \
    -servername "${HOST}" \
    </dev/null 2>/dev/null \
    | openssl x509 -noout -enddate 2>/dev/null \
    | cut -d= -f2)

  if [[ -z "$EXPIRY" ]]; then
    echo "ERROR: could not retrieve certificate for ${HOST_PORT}" >&2
    STATUS=1
    continue
  fi

  EXPIRY_EPOCH=$(date -d "${EXPIRY}" +%s)
  NOW_EPOCH=$(date +%s)
  DAYS_LEFT=$(( (EXPIRY_EPOCH - NOW_EPOCH) / 86400 ))

  if (( DAYS_LEFT <= WARN_DAYS )); then
    echo "WARN: ${HOST_PORT} expires in ${DAYS_LEFT} days (${EXPIRY})"
    STATUS=1
  else
    echo "OK:   ${HOST_PORT} expires in ${DAYS_LEFT} days (${EXPIRY})"
  fi
done

exit $STATUS
```

Run it as:

```bash
chmod +x check-cert-expiry.sh
WARN_DAYS=14 ./check-cert-expiry.sh example.com:443 smtp.example.com:587
```

### 5. Common Error Codes and What They Mean

| Error / Alert | Meaning | Typical Cause |
|---|---|---|
| `verify error:num=2` | Unable to get issuer certificate | Intermediate CA missing from chain |
| `verify error:num=10` | Certificate has expired | Cert past `Not After` date |
| `verify error:num=18` | Self-signed certificate | No trusted CA for self-signed cert |
| `verify error:num=20` | Unable to get local issuer | Root CA not in system trust store |
| `alert handshake failure (40)` | Handshake failed | No common cipher or protocol version |
| `alert unknown ca (48)` | Unknown certificate authority | Client does not trust the signing CA |
| `alert certificate expired (45)` | Certificate expired | Server cert is past its validity window |

For deeper inspection, set `OPENSSL_DEBUG=1` or compile with `--enable-ktls` and trace the TLS record layer. For most production issues, the verify error number and the certificate chain output from `s_client -showcerts` are sufficient to find the problem.
