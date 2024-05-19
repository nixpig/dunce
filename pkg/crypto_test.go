package pkg

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

var mockCrypto = new(MockCrypto)

func TestCrypto(t *testing.T) {
	scenarios := map[string]func(t *testing.T, crypto Crypto){
		"generate password (success)":         testGenerateFromPasswordSuccess,
		"generate password (error)":           testGenerateFromPasswordError,
		"compare hash and password (success)": testCompareHashAndPasswordSuccess,
		"compare hash and password (error)":   testCompareHashAndPasswordError,
	}

	for scenario, fn := range scenarios {
		crypto := NewCryptoImpl(
			mockCrypto.passwordGenerator,
			mockCrypto.hashAndPasswordComparer,
		)

		t.Run(scenario, func(t *testing.T) {
			fn(t, crypto)
		})
	}
}

type MockCrypto struct {
	mock.Mock
}

func (m *MockCrypto) passwordGenerator(
	password []byte,
	cost int,
) ([]byte, error) {
	args := m.Called(password, cost)

	return args.Get(0).([]byte), args.Error(1)
}

func (m *MockCrypto) hashAndPasswordComparer(
	hashedPassword []byte,
	password []byte,
) error {
	args := m.Called(hashedPassword, password)

	return args.Error(0)
}

func testGenerateFromPasswordSuccess(t *testing.T, crypto Crypto) {
	mockCryptoPasswordGenerator := mockCrypto.On("passwordGenerator", []byte("p4ssw0rd"), 12).
		Return([]byte("generated"), nil)

	password, err := crypto.GenerateFromPassword([]byte("p4ssw0rd"), 12)

	require.NoError(t, err, "should not error")
	require.Equal(
		t,
		[]byte("generated"),
		password,
		"should return value from generator",
	)

	mockCryptoPasswordGenerator.Unset()
	if res := mockCrypto.AssertExpectations(t); !res {
		t.Error("should call the password generator")
	}
}

func testGenerateFromPasswordError(t *testing.T, crypto Crypto) {
	mockCryptoPasswordGenerator := mockCrypto.On("passwordGenerator", []byte("p4ssw0rd"), 12).
		Return([]byte(""), errors.New("generator_error"))

	password, err := crypto.GenerateFromPassword([]byte("p4ssw0rd"), 12)

	require.Empty(t, password, "shouldn't return a generated value")
	require.EqualError(
		t,
		err,
		"generator_error",
		"should return error from generator",
	)

	mockCryptoPasswordGenerator.Unset()
	if res := mockCrypto.AssertExpectations(t); !res {
		t.Error("should call password generator")
	}
}

func testCompareHashAndPasswordSuccess(t *testing.T, crypto Crypto) {
	mockCryptoHashAndPasswordComparer := mockCrypto.On("hashAndPasswordComparer", []byte("h4shedp4ssw0rd"), []byte("p4ssw0rd")).
		Return(nil)

	err := crypto.CompareHashAndPassword(
		[]byte("h4shedp4ssw0rd"),
		[]byte("p4ssw0rd"),
	)

	require.NoError(t, err, "should not return error")

	mockCryptoHashAndPasswordComparer.Unset()
	if res := mockCrypto.AssertExpectations(t); !res {
		t.Error("should call hash and password comparer")
	}
}

func testCompareHashAndPasswordError(t *testing.T, crypto Crypto) {
	mockCryptoHashAndPasswordComparer := mockCrypto.On("hashAndPasswordComparer", []byte("h4shedp4ssw0rd"), []byte("p4ssw0rd")).
		Return(errors.New("comparer_error"))

	err := crypto.CompareHashAndPassword(
		[]byte("h4shedp4ssw0rd"),
		[]byte("p4ssw0rd"),
	)

	require.EqualError(t, err, "comparer_error", "should return error")

	mockCryptoHashAndPasswordComparer.Unset()
	if res := mockCrypto.AssertExpectations(t); !res {
		t.Error("should call hash and password comparer")
	}
}
