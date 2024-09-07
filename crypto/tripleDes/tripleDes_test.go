package tripleDes_test

import (
	"fmt"
	"github.com/odycenter/std-library/crypto"
	"github.com/odycenter/std-library/crypto/tripleDes"
	"github.com/odycenter/std-library/json"
	"testing"
)

func TestTripledes_Encrypt(t *testing.T) {
	td := tripleDes.New("aaaaaaaaaaaaaaaaaaaaaaaa", "bbbbbbbbbbbbbbbbbbbbbbbb", crypto.PaddingPKCS7, crypto.ECB)
	fmt.Println(td.Encrypt("欢迎使用library ").Hex())
}

func TestTripledes_Decrypt(t *testing.T) {
	td := tripleDes.New("aaaaaaaaaaaaaaaaaaaaaaaa", "bbbbbbbbbbbbbbbbbbbbbbbb", crypto.PaddingPKCS7, crypto.ECB).WithHex()
	fmt.Println(td.Decrypt([]byte("9c21ec39e3a5a3fec3763c394817b599fc1caf0d699ad364")).String())
}

type T struct {
	Type                    string `json:"type"`
	ProjectId               string `json:"project_id"`
	PrivateKeyId            string `json:"private_key_id"`
	PrivateKey              string `json:"private_key"`
	ClientEmail             string `json:"client_email"`
	ClientId                string `json:"client_id"`
	AuthUri                 string `json:"auth_uri"`
	TokenUri                string `json:"token_uri"`
	AuthProviderX509CertUrl string `json:"auth_provider_x509_cert_url"`
	ClientX509CertUrl       string `json:"client_x509_cert_url"`
}

func TestAll(t *testing.T) {
	v := T{
		Type:                    "service_account",
		ProjectId:               "airy-box-379304",
		PrivateKeyId:            "b148efb1a5faa789f0df09c540dd6b379a248bc8",
		PrivateKey:              "-----BEGIN PRIVATE KEY-----\nMIIEvgIBADANBgkqhkiG9w0BAQEFAASCBKgwggSkAgEAAoIBAQCzY2dGc38P0PqQ\n9TfL+yEi2FywjogcF/PblSKEY4rTnZ7Mc6++SnL60vbxltiuS0gQm7xjgWSbVAXL\nPPB+8xtp4ykSanM2qQ1kSveIFNYQMTZ/N+D/3fDSNJ8AS6/LvxBepLuQPrw2QZug\nTEVVAd+hHLQXbKPtkcRzsQkhoMY+ISqra5FHZs47Wmt7btxMY92NrTf44B8folhg\nUkkmzNF6k2WW5ZX0RDQJ64DT0uB3zPUGhS4I44en8TVWJQEueAmvdKOX7Er+jU/P\neCgIhUnKpH0t836nH55NysqJm9bLsETe0/PYiJb0ANu3jkjabPvqzzlZWvG8bE1T\nQJK5nuMrAgMBAAECggEAF8HTDD+A1QVmHwslZ+b1onm+UhuY5wTnuhhBdABYLbey\n2id6tNwt52NHv62skJijuJ9cubOJ3BBVRsSRTR7BGKre5kjMGAcXdkV/YWRGXT/x\n6bNAIbgGNrCFW4f/23SNpHBzQ/fVppHVxyA9Zl6oe9vN9fQ2BDZz1VEy1mNGczUU\ncInbcAvCXMpmI9rE72vqnN37S7GTMf3G68CGybiCZLfLlhkg1c7V+XCkJg60MAZ1\nmxYXEGr6M/OoxRza9LBFgMw97UMiKEKtLQ+99ObwmuMQbJ6g/YU5WiTw+5LlbvGc\ngMmRzjM+PjM0Nsz3i5aI15a68d3PmeBPUOF8wu16UQKBgQDgoKrYLsaZZiZZvQnn\nKdpEhTBRrYMI9n9Rd2XgrsTt/bqQ/HJpoKBzmZ5Fh5pASXrPiKFPTRrQUjDE8UYh\n5CBK6PuaI3ZJr7ShwKAnQhyJZS+ib4MxtWedlCM2riQtdfrJA5VSr78sHYYjlT6B\ny5HYu8ExZjFe2enYDxBs2KHOfQKBgQDMcUAQpCmadSThmY6uFbYddr5f3TjSBQF+\nTSAqlJBPo/PHvxnzoPJkiu17cHmvaxn/iOkirQP2vgMTxofq++hskddIX5EPRHB/\n0SbX9EICWmA/aTmbpzAA5rH3TOj+AsbIztCsS7KxXieW0uUct0CgyoIHxmLFaOSf\nbeP2UfLgxwKBgC8hpxc7IVKYe12C66wEPRb5dzz8Ei10Qxyd19N/+DQTc+zt+zes\ni14WEn52SGhKwqj++xG/lOu3AyKfmV6NFjWBkyExZaVqZ5U07KWwGnq9r3P+v+FT\nNc17grP7b/3V7mv1A4TY+VzRSQ74RqhHRW/bXVr3HU7QnF9IMeMUxUalAoGBAMNW\nKdo8oCueZgDQEY2v3PPF8xvxaUrx0X11/5fvnvsZMeHWa9tmGnOKcmIRE5NSB+Mq\nU2b4XOMypgoNFOymiGrD5iiWdylZQQw/MJgCH9fTtkagKZTZZ3pU8hHSAIRC7uAL\nC0K0iSYDSlxHYPXQ+gUnuJnpKZJpKJhUDQ3bOu8dAoGBAM/sz4VD+Ws43y3P1vZ8\n3/VgNP5Q5z47PCLDzFS77zBGBXvrunCnHCiMz2huVhbZ11XuGDx/ndF2SKZF3Bv5\npANHwYkDnD255G6D7v6s8qOjDpoYMOxAmGBTTNtAqvd5BIEdUYHRYEc+1h1dThWt\nq25L/HbMqJqxjcsxTlOnWLOq\n-----END PRIVATE KEY-----\n",
		ClientEmail:             "ocr-680@airy-box-379304.iam.gserviceaccount.com",
		ClientId:                "107076889875450634572",
		AuthUri:                 "https://accounts.google.com/o/oauth2/auth",
		TokenUri:                "https://oauth2.googleapis.com/token",
		AuthProviderX509CertUrl: "https://www.googleapis.com/oauth2/v1/certs",
		ClientX509CertUrl:       "https://www.googleapis.com/robot/v1/metadata/x509/ocr-680%40airy-box-379304.iam.gserviceaccount.com",
	}
	cipher := tripleDes.New("eShVmYp3s6v9y$B&E)H@McQf", "jWnZr4u7", crypto.PaddingPKCS7, crypto.CFB)
	fmt.Println(cipher.Encrypt(string(json.Stringify(v))).Hex())
}
