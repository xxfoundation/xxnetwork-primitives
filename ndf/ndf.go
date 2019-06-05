package globals

import (
	"encoding/base64"
	"encoding/json"
	"github.com/spacemonkeygo/openssl"
	"strings"
)

func CreateNDFPublicKey() openssl.PublicKey {
	publicKey, err := openssl.LoadPublicKeyFromPEM(
		[]byte(`-----BEGIN PUBLIC KEY-----
MIIDRjCCAjkGByqGSM44BAEwggIsAoIBAQDrFa5KY+P5MAv9C3cc5ndlLzU2dnjk
5eFln7/t9OwGhNs63UcZW2iYELnM30wVx7gYQigPUyXP6973w3R9k7/6QPXZ20IP
SDFr3Q5OstIcCgOPUTMHycYkRbLQzZP+BHaTE+RtOC1uMKQYF/7rDLvkqia1+3p2
e7OSWe6HV2k+lbaj4DHTHOtFqSe6QjEjVobNkF0kR0Y/kFIl8lYL8+8V2MMKVlb6
8YN0T2LmZXPCWwB0nDeCQ4ryvFMlXoS6x497PGgHEk3qQtnVAc02q9HQl0r9i0sI
yNburnAbmvm3jS76sTr1qOAGi3GaaY2nSDe0c0OAcTNXwcldx6SuicBXAiEAr+XF
Puot+mcFb+sdZldhFKbIRehp0u/M0wYZd+lywb0CggEAFtHMikAwFcLFtxZIDA3b
DIWRUN1zWJ7xVL19lnPxnxVSbBXBYbXK5PKRN/N3eGcVF6pq4N8cqOAco/2lntSS
YaATNSVieJGlVsN9OOhztFlEahN97bqucbu7ikYKitC1hHFQzapeejw1pjMDlxOt
bxSfn0X4kujWkOSI0AGj7CpcTilp0ntnzDSEEu2R8Ta+7L2VXsb3fNxxkWWfc+xo
X2DsPBCa+zSmwIXuRsEethH7MAyzf3uPA98ZZRG7eGZfW1+/fgl5GXlIUGbxzf0K
FoA0Qn2QDZ+8HodT147JsyxHyKnSEhs5wviq9b8GhsxEEw8DCKt8GU40q6WQCaUj
vAOCAQUAAoIBAGkgeVn5Dg5x7Hn0W9dPlVNJXzJEsJGZItu+/VrzGA3ZVfdeb9/e
579cfA7DuoHpZnLA8+Bb8kMuaMlzaNA0Vo6tfJsqmZKq+DH+Ww7Jhs6OJvdlG3ik
0NktvkN+MiYke301xbdZcbRtWu97e2930onIMZt+QnAq1+7UQt81Eav4bFtOZyR5
xDCxFHzFn2/wNlg/m/QWgAFQaRHAumc2X+26qUyblXS10AxoPEIWXs0DGgY1MQzk
7N/zhUPhHmEBTAhFiXqnhV4v2IFAGf4GPwmF9kIrQkQ3paI2mh8EW1qYEl7BD+jw
fcIgcc/yqy1i6kUTYCDRHz1h7+vTHcku1TA=
-----END PUBLIC KEY-----`))
	if err != nil {
		panic(err)
	}
	return publicKey
}

// This struct is currently generated in Terraform and decoded here
// So, if the way it's generated in Terraform changes, we also need to change
// the struct
// TODO Use UnmarshalJSON for user and node IDs and groups, at the least
//  We also need to unmarshal the Timestamp to a time.Time
//  See https://blog.gopheracademy.com/advent-2016/advanced-encoding-decoding/
//  for information about how to do this.
type NetworkDefinitionJSON struct {
	Timestamp string
	Gateways []GatewayInfoJSON
	Nodes []NodeInfoJSON
	Registration RegistrationInfoJSON
	Udb UDBInfoJSON
	E2e GroupJSON
	Cmix GroupJSON
}

type GatewayInfoJSON struct{
	Address string
	Tls_certificate string
}

type NodeInfoJSON struct {
	Id []byte
	Dsa_public_key string
	Address string
	Tls_certificate string
}

type RegistrationInfoJSON struct {
	Dsa_public_key string
    Address string
}

type UDBInfoJSON struct {
	Id []byte
	Dsa_public_key string
}

type GroupJSON struct {
	Prime string
	Small_prime string
	Generator string
}

// Returns an error if base64 signature decodes incorrectly
// Returns an error if signature verification fails
// Otherwise, returns an object from the json with the contents of the file
func DecodeNDF(ndf string) (*NetworkDefinitionJSON, error) {
	var jsonString string
    jsonString, err := verify(ndf)
	if err != nil {
		return nil, err
	}
    networkDefinition := &NetworkDefinitionJSON{}
	err = json.Unmarshal([]byte(jsonString), networkDefinition)
	if err != nil {
        return nil, err
	}
	return networkDefinition, err
}

func verify(ndf string) (jsonString string, err error) {
	lines := strings.Split(ndf, "\n")
	// Base64 decode the signature to a byte slice
	// VerifyPKCS1v15 requires a raw byte slice to do its processing
	signature, err := base64.StdEncoding.DecodeString(lines[1])
	if err != nil {
		return "", err
	}
	publicKey := CreateNDFPublicKey()
	return lines[0],
		publicKey.VerifyPKCS1v15(openssl.SHA256_Method, []byte(lines[0]), signature)
}
