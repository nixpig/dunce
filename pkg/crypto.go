package pkg

import "golang.org/x/crypto/bcrypt"

type Crypto interface {
	GenerateFromPassword(password []byte, cost int) ([]byte, error)
	CompareHashAndPassword(hashedPassword []byte, password []byte) error
}

type CryptoImpl struct{}

func NewCryptoImpl() CryptoImpl {
	return CryptoImpl{}
}

func (c CryptoImpl) GenerateFromPassword(password []byte, cost int) ([]byte, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword(password, cost)
	if err != nil {
		return nil, err
	}

	return hashedPassword, nil
}

func (c CryptoImpl) CompareHashAndPassword(hashedPassword []byte, password []byte) error {
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)); err != nil {
		return err
	}

	return nil
}
