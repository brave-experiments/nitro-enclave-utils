package certificate_test

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"regexp"
	"testing"

	"github.com/brave/nitriding/certificate"
	"github.com/brave/nitriding/nitridingtest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMakeBasePrivilegedCert_HappyPath(t *testing.T) {
	cert, err := certificate.MakeBasePrivilegedCert("", "")
	assert.NoError(t, err)
	assert.NotNil(t, cert)

	VerifyCert(t, cert)

	parsedCert, err := x509.ParseCertificate(cert.DerBytes())
	assert.NoError(t, err)
	privateKey := cert.PrivateKey()
	require.NoError(t, err)
	assert.Equal(t, parsedCert.PublicKey, &privateKey.PublicKey)
}

func TestMakeBasePrivilegedCert_Errors(t *testing.T) {
	tests := map[string]struct {
		certOrg  string
		fqdn     string
		errRegex string
	}{
		"bad CertOrg": {
			certOrg:  "\x92",
			fqdn:     "",
			errRegex: "asn1: string not valid UTF-8",
		},
		"bad FQDN": {
			certOrg:  "",
			fqdn:     "\xff",
			errRegex: "x509: .* cannot be encoded as an IA5String",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			cert, err := certificate.MakeBasePrivilegedCert(tc.certOrg, tc.fqdn)
			assert.Error(t, err)
			assert.Regexp(t, regexp.MustCompile(tc.errRegex), err.Error())
			assert.Equal(t, certificate.BasePrivilegedCert{}, cert)
		})
	}
}

func FuzzMakeBasePrivilegedCert(f *testing.F) {
	tests := []struct {
		certOrg string
		fqdn    string
	}{
		{"", ""},     // happy path
		{"\x92", ""}, // asn1: string not valid UTF-8
		{"", "\xff"}, // x509: "\xff" cannot be encoded as an IA5String
	}
	for _, tc := range tests {
		f.Add(tc.certOrg, tc.fqdn)
	}
	f.Fuzz(func(t *testing.T, certOrg string, fqdn string) {
		cert, makeErr := certificate.MakeBasePrivilegedCert(certOrg, fqdn)
		if makeErr == nil {
			VerifyCert(t, cert)
		} else {
			assert.Equal(t, certificate.BasePrivilegedCert{}, cert)

			ok, err := nitridingtest.ErrorMatchesPattern(
				makeErr,
				"asn1: string not valid UTF-8",
				"x509: .* cannot be encoded as an IA5String",
			)
			assert.NoError(t, err)
			assert.Truef(t, ok, "CertOrg '%v', FQDN '%#v', error '%v'", certOrg, fqdn, makeErr)
		}
	})
}

func TestBasePrivilegedCert_PrivateKey_HappyPath(t *testing.T) {
	derBytes := certificate.DerBytes("some DER bytes")

	privateKey, err := ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	assert.NoError(t, err)

	cert, err := certificate.MakeBasePrivilegedCertFromRaw(
		derBytes,
		certificate.BaseConverter{},
		privateKey,
	)
	assert.NoError(t, err)

	outPrivateKey := cert.PrivateKey()
	assert.NoError(t, err)
	assert.Equal(t, *privateKey, *outPrivateKey)
}

func TestMakeTLSCertificate_HappyPath(t *testing.T) {
	cert, err := certificate.MakeBasePrivilegedCert("", "")
	assert.NoError(t, err)
	VerifyCert(t, cert)

	tlsCert, err := cert.TLSCertificate()
	assert.NoError(t, err)
	assert.Equal(t, cert.PrivateKey(), tlsCert.PrivateKey)
}

func TestBasePrivilegedCert_TLSCertificate_NilCertPrivateKey(t *testing.T) {
	cert, err := certificate.MakeBasePrivilegedCertFromRaw(
		certificate.DerBytes{},
		certificate.BaseConverter{},
		nil,
	)
	assert.NoError(t, err)

	tlsCert, err := cert.TLSCertificate()
	assert.ErrorContains(t, err, "privateKey cannot be nil")
	assert.Equal(t, tls.Certificate{}, tlsCert)
}

func TestBasePrivilegedCert_TLSCertificate_BadCertPrivateKey(t *testing.T) {
	cert, err := certificate.MakeBasePrivilegedCertFromRaw(
		certificate.DerBytes{},
		certificate.BaseConverter{},
		&ecdsa.PrivateKey{},
	)
	assert.NoError(t, err)

	tlsCert, err := cert.TLSCertificate()
	assert.ErrorContains(t, err, "x509: unknown curve while marshaling")
	assert.Equal(t, tls.Certificate{}, tlsCert)
}
