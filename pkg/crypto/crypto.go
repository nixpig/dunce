package crypto

type PasswordGenerator func(password []byte, cost int) ([]byte, error)
type HashAndPasswordComparer func(hashedPassword []byte, password []byte) error

type Crypto interface {
	GenerateFromPassword(password []byte, cost int) ([]byte, error)
	CompareHashAndPassword(hashedPassword []byte, password []byte) error
}

type CryptoImpl struct {
	generateFromPassword   PasswordGenerator
	compareHashAndPassword HashAndPasswordComparer
}

func NewCryptoImpl(
	generateFromPassword PasswordGenerator,
	compareHashAndPassword HashAndPasswordComparer,
) CryptoImpl {
	return CryptoImpl{
		generateFromPassword:   generateFromPassword,
		compareHashAndPassword: compareHashAndPassword,
	}
}

func (c CryptoImpl) GenerateFromPassword(
	password []byte,
	cost int,
) ([]byte, error) {
	return c.generateFromPassword(password, cost)
}

func (c CryptoImpl) CompareHashAndPassword(
	hashedPassword []byte,
	password []byte,
) error {
	return c.compareHashAndPassword(hashedPassword, password)
}
