package services

import "os"
import "fmt"
import "errors"
import "io/ioutil"
import "crypto/rsa"
import "crypto/x509"
import "encoding/pem"
import "github.com/jinzhu/gorm"
import "github.com/SermoDigital/jose/jws"
import "github.com/SermoDigital/jose/crypto"

import "github.com/dadleyy/charcoal.api/models"

type UserClientManager struct {
	*gorm.DB
}

func decodePrivateJwtToken(path string) (*pem.Block, error) {
	buf, err := ioutil.ReadFile(path)

	if err != nil {
		return nil, err
	}

	block, rest := pem.Decode(buf)

	if len(rest) > 0 {
		return nil, fmt.Errorf("BAD_PEM_DECODE")
	}

	return block, nil
}

func (engine *UserClientManager) Validate(input string, client *models.Client) error {
	var token models.ClientToken
	var user models.User

	where := engine.Where("token = ?", input)

	if err := where.First(&token).Error; err != nil {
		return errors.New("ERR_BAD_TOKEN")
	}

	if token.Client != client.ID {
		return errors.New("ERR_BAD_TARGET_CLIENT")
	}

	if err := engine.First(&user).Error; err != nil {
		return errors.New("ERR_NO_USER_FOR_CLIENT")
	}

	block, err := decodePrivateJwtToken(os.Getenv("JWT_PUBLIC_KEY"))

	if err != nil {
		return errors.New("ERR_BAD_PUBLIC_KEY")
	}

	rsapub, err := x509.ParsePKIXPublicKey(block.Bytes)

	if err != nil {
		return errors.New("ERR_RSA_PUBLIC_KEY_PARSE")
	}

	buffer := []byte(token.Token)
	parsed, err := jws.ParseJWT(buffer)

	if err != nil {
		return errors.New("ERR_BAD_JWT_PARSE")
	}

	if err := parsed.Validate(rsapub, crypto.SigningMethodRS512); err != nil {
		return errors.New("ERR_BAD_VALIDATE")
	}

	return nil
}

func (engine *UserClientManager) Associate(user *models.User, client *models.Client) (models.ClientToken, error) {
	var result models.ClientToken
	var rsapriv *rsa.PrivateKey

	if client == nil || user == nil {
		return models.ClientToken{}, fmt.Errorf("INVALID_RECORDS")
	}

	if client.ID == 0 {
		return models.ClientToken{}, fmt.Errorf("BAD_CLIENT_ID")
	}

	result = models.ClientToken{
		Client: client.ID,
		User:   user.ID,
	}

	var tcount uint

	if err := engine.Model(&result).Where(result).First(&result).Count(&tcount).Error; err != nil && tcount != 0 {
		return result, fmt.Errorf("FAILED_COUNT: %s", err.Error())
	}

	if tcount >= 1 {
		return result, nil
	}

	block, err := decodePrivateJwtToken(os.Getenv("JWT_PRIVATE_KEY"))

	if err != nil {
		return models.ClientToken{}, fmt.Errorf("BAD_PRIVATE_KEY_DECODE: %s", err.Error())
	}

	if x509.IsEncryptedPEMBlock(block) {
		return models.ClientToken{}, fmt.Errorf("ENCRYPTED_BLOCK")
	}

	rsapriv, err = x509.ParsePKCS1PrivateKey(block.Bytes)

	if err != nil {
		return models.ClientToken{}, fmt.Errorf("BAD_PARSE_BLOCK: %s", err.Error())
	}

	claims := jws.Claims{
		"token": fmt.Sprintf("%s:%d", client.ClientID, user.ID),
	}

	token := jws.NewJWT(claims, crypto.SigningMethodRS512)

	sbuf, err := token.Serialize(rsapriv)

	if err != nil {
		return models.ClientToken{}, fmt.Errorf("BAD_SERIALIZATION: %s", err.Error())
	}

	result.Token = string(sbuf)

	if err := engine.Save(&result).Error; err != nil {
		return models.ClientToken{}, err
	}

	return result, nil
}
