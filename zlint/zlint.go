/* zlint.go
 * Used to check parsed info from certificate for compliance
 */

package zlint

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/zmap/zcrypto/x509"
	"github.com/zmap/zlint/lints"
)

type PrettyOutput struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Provenance  string `json:"provenance"`
}

//Pretty Print lint outputs
func PrettyPrintZLint() {
	for _, l := range lints.Lints {
		p := PrettyOutput{}
		p.Name = l.Name
		p.Description = l.Description
		p.Provenance = l.Provenance

		buffer := new(bytes.Buffer)
		enc := json.NewEncoder(buffer)
		enc.SetEscapeHTML(false)
		enc.Encode(p)
		fmt.Print(buffer.String())
	}
}

//Calls all other checks on parsed certs producing version
func ZLintResultTestHandler(cert *x509.Certificate) *lints.ZLintResult {
	if cert == nil {
		return nil
	}
	//run all tests
	ZLintResult := lints.ZLintResult{}
	ZLintResult.Execute(cert)
	ZLintResult.ZLintVersion = lints.ZLintVersion
	ZLintResult.Timestamp = time.Now().Unix()
	return &ZLintResult
}

//Calls all other checks on parsed certs
func ParsedTestHandler(cert *x509.Certificate) (map[string]string, error) {
	if cert == nil {
		return nil, errors.New("zlint: nil pointer passed in, no data returned")
	}
	//run all tests
	var out map[string]string = make(map[string]string)
	for _, l := range lints.Lints {
		result, err := l.ExecuteTest(cert)
		if err != nil {
			return out, err
		}
		out[l.Name] = lints.EnumToString(result.Result)
	}

	return out, nil
}

//expects Base64 encoded ASN.1 DER string, wrapper for ParsedTestHandler
func Lint64(certIn string) (map[string]string, error) {

	der, err := base64.StdEncoding.DecodeString(certIn) //decode string
	if err != nil {
		return nil, err //parsing error, no tests performed
	}

	cert, err := x509.ParseCertificate(der)
	if err != nil {
		return nil, err //parsing error, no tests performed
	}

	return ParsedTestHandler(cert) //return available reports & error from main testing function
}
