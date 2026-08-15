---
title: "Keyless Container Image Signing with Cosign and OIDC"
slug: b4c9a1d2
aliases: [cosign-keyless-container-signing]
date: 2026-04-12
tags: [security, devops, docker]
lang: en
draft: false
type: post
---

Keyless container image signing removes the operational burden of managing long-lived cryptographic keys. Sigstore pairs OpenID Connect (OIDC) tokens with short-lived X.509 certificates and an append-only transparency log. This setup provides verifiable container integrity in continuous integration pipelines and Kubernetes runtime environments.

## Fun Facts

**Fact 1.** Fulcio issues ephemeral certificates valid for only ten minutes, binding identity tokens from GitHub Actions or cloud providers to single-use public keys generated in memory.

**Fact 2.** Rekor records every signature in a Merkle tree transparency log based on Google Trillian, making signatures publicly auditable and resistant to retroactive modification.

**Fact 3.** Sigstore became a Linux Foundation project in 2021 and serves as the primary artifact signing mechanism for Kubernetes official releases since version 1.24.

---

## Tips and Tricks

### 1. Enable GitHub Actions OIDC Token Permissions
Grant workflow jobs access to request short-lived OIDC tokens directly from the GitHub token service.

```yaml
name: build-and-sign
on:
  push:
    tags:
      - 'v*'

permissions:
  contents: read
  id-token: write
  packages: write
```

### 2. Sign Container Images Without Key Files
Run `cosign sign` in your workflow. Cosign detects the ambient GitHub OIDC token, requests a certificate from Fulcio, and logs the payload to Rekor.

```bash
cosign sign --yes \
  ghcr.io/example/app:v1.2.0@sha256:e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
```

### 3. Verify Container Signatures on the Command Line
Verify that an image was signed by a specific workflow repository and identity provider.

```bash
cosign verify \
  --certificate-identity-regexp="https://github.com/example/app/\.github/workflows/.*" \
  --certificate-oidc-issuer="https://token.actions.githubusercontent.com" \
  ghcr.io/example/app:v1.2.0
```

### 4. Enforce Image Signatures in Kubernetes with Kyverno
Deploy a Kyverno policy to reject unsigned pod images or images signed outside approved workflows.

```yaml
apiVersion: kyverno.io/v1
kind: ClusterPolicy
metadata:
  name: enforce-keyless-signatures
spec:
  validationFailureAction: Enforce
  background: false
  rules:
    - name: check-image-signature
      match:
        any:
          - resources:
              kinds:
                - Pod
      verifyImages:
        - imageReferences:
            - "ghcr.io/example/*"
          attestors:
            - entries:
                - keyless:
                    issuer: "https://token.actions.githubusercontent.com"
                    subjectRegexp: "https://github.com/example/app/\.github/workflows/.*"
                    rekor:
                      url: "https://rekor.sigstore.dev"
```

### 5. Query Rekor Transparency Log Entries
Search Rekor for signed entry metadata by container digest to inspect timestamps and verification certificates.

```bash
rekor-cli search \
  --artifact ghcr.io/example/app@sha256:e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
```
