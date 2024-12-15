package utils_test

import (
	"testing"

	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/utils"
	"github.com/stretchr/testify/require"
)

func TestCreate(t *testing.T) {
	secret := "this_is_16_byteK"
	ct, err := utils.NewCryptToken(secret)
	require.NoError(t, err)

	var userID uint64 = uint64(12)
	userType := "applicant"
	SessionID := "session12e2325wetr4"
	
	cypher, err := ct.Create(userID, userType, SessionID)
	require.NoError(t, err)

	ok, err := ct.Check(userID, userType, SessionID, cypher)
	require.NoError(t, err)
	require.True(t, ok)
}
