package utils

import (
	"encoding/base64"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncryptAES128ECB(t *testing.T) {
	enc, err := EncryptAES128ECB([]byte(`{"document_id":"IdDoc_002","status_document":"complete","result":"00","email_user":"testing5@digisign.id","notif":"Sukses"}`), []byte("RBazsYSDTuShYbUG"))
	assert.NoError(t, err)
	res := base64.StdEncoding.EncodeToString(enc)
	esc := url.QueryEscape(res)

	uesc, err := url.QueryUnescape(esc)
	assert.NoError(t, err)
	dec, err := base64.StdEncoding.DecodeString(uesc)
	assert.NoError(t, err)
	ress := DecryptAES128ECB([]byte(dec), []byte("RBazsYSDTuShYbUG"))
	assert.Equal(t, `{"document_id":"IdDoc_002","status_document":"complete","result":"00","email_user":"testing5@digisign.id","notif":"Sukses"}`, string(ress))
}

func TestDecryptAES128ECB(t *testing.T) {
	esc, err := url.QueryUnescape("Cij6Jz9ui76E5Ky%2FqJRqt5jMp1gx%2Bz1onr2%2FO%2Fb%2FQQOlCOI%2FsNbocobyzf4XkeXTCZRJ%2B7r7J3%2FI%0AjMKCjBpOTpZuuhUUeQTD%2BmU43vsKbSKj0oQxhKUvUG%2BNcdYcTSRiM3ikRaMiBB4rrxQXdEwUu18r%0A1naOnZwvJw0ZnEx94Z4%3D")
	assert.NoError(t, err)

	dec, err := base64.StdEncoding.DecodeString(esc)
	assert.NoError(t, err)
	res := DecryptAES128ECB([]byte(dec), []byte("RBazsYSDTuShYbUG"))
	assert.Equal(t, `{"document_id":"IdDoc_002","status_document":"complete","result":"00","email_user":"testing5@digisign.id","notif":"Sukses"}`, string(res))
}
