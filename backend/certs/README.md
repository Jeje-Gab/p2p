# SSL/TLS Certificates

This directory contains SSL/TLS certificates for HTTPS support.

## Auto-Generated Certificates

The certificates in this directory are **self-signed** and generated automatically for **development purposes only**.

### Files

- `server.key` - Private key (4096-bit RSA)
- `server.crt` - Public certificate (valid for 365 days)
- `openssl.cnf` - OpenSSL configuration file

## Regenerating Certificates

If you need to regenerate the certificates (e.g., they expired), run:

```bash
cd backend/certs
openssl req -x509 -newkey rsa:4096 -keyout server.key -out server.crt -days 365 -nodes -config openssl.cnf
```

## Certificate Details

- **Common Name (CN):** localhost
- **Organization:** CS2 P2P Skins
- **Validity:** 365 days from generation
- **Key Size:** 4096-bit RSA
- **Subject Alternative Names:**
  - DNS: localhost, *.localhost
  - IP: 127.0.0.1, ::1

## Production Deployment

⚠️ **IMPORTANT:** Self-signed certificates are **NOT suitable for production**.

For production, use certificates from a trusted Certificate Authority (CA):

### Option 1: Let's Encrypt (Free)
```bash
# Using Certbot
sudo certbot certonly --standalone -d yourdomain.com
```

### Option 2: Commercial CA
Purchase certificates from:
- DigiCert
- GlobalSign
- Sectigo
- Others

### Option 3: Cloud Provider Certificates
- AWS Certificate Manager (ACM)
- Google Cloud Certificate Manager
- Azure Key Vault Certificates

## Security Notes

1. **Never commit private keys to version control**
   - `*.key` files are in `.gitignore`

2. **Rotate certificates regularly**
   - Self-signed: regenerate every 365 days
   - Production: follow CA renewal schedule

3. **Protect private keys**
   - Set proper file permissions: `chmod 600 server.key`
   - Store in secure location
   - Use hardware security modules (HSM) in production

4. **Use strong ciphers**
   - TLS 1.2+ only
   - Disable weak ciphers (see server configuration)

## Trusting Self-Signed Certificates (Development)

### For Bruno/API Testing
In Bruno settings:
- Settings → SSL/TLS → Disable SSL verification ✓

### For Browsers
Add exception for `https://localhost:8443`

### For System-Wide Trust (Optional)

**Windows:**
```powershell
# Import certificate to Trusted Root CA
Import-Certificate -FilePath server.crt -CertStoreLocation Cert:\LocalMachine\Root
```

**Linux:**
```bash
sudo cp server.crt /usr/local/share/ca-certificates/cs2-p2p.crt
sudo update-ca-certificates
```

**macOS:**
```bash
sudo security add-trusted-cert -d -r trustRoot -k /Library/Keychains/System.keychain server.crt
```

## Verification

Verify certificate details:
```bash
openssl x509 -in server.crt -text -noout
```

Test HTTPS connection:
```bash
curl -v --cacert server.crt https://localhost:8443/healthz
# or ignore cert verification
curl -k https://localhost:8443/healthz
```

---

**Generated for:** CS2 P2P Skins Trading Platform
**Purpose:** Development and Testing
**Not for Production Use**
