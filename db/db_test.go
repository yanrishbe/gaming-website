package db

import (
	"testing"

	"github.com/yanrishbe/gaming-website/entities"

	"github.com/stretchr/testify/require"
)

type testUsers struct {
	requestUser *entities.User
	resultUser  *entities.User
}

var tableUsers = []testUsers{
	{
		&entities.User{Balance: 400},
		&entities.User{ID: 1, Balance: 100},
	},
	{
		&entities.User{Balance: 0},
		&entities.User{Balance: 300},
	},
	{
		&entities.User{Balance: 200},
		&entities.User{ID: 2, Balance: 100},
	},
}

var mockMap = New()

func TestDB_SaveUser(t *testing.T) {
	require := require.New(t)
	initialLen := len(mockMap.UsersMap)
	errSave := mockMap.SaveUser(tableUsers[0].requestUser)
	require.NoError(errSave)
	require.Equal(tableUsers[0].resultUser.Balance, tableUsers[0].requestUser.Balance)
	require.Equal(initialLen+1, len(mockMap.UsersMap))
}

func TestDB_DeleteUser(t *testing.T) {
	require := require.New(t)
	mockMap.UsersMap[tableUsers[0].resultUser.ID] = tableUsers[0].resultUser
	initialLen := len(mockMap.UsersMap)
	require.NoError(mockMap.DeleteUser(tableUsers[0].resultUser.ID))
	require.Equal(initialLen-1, len(mockMap.UsersMap))
}

func TestDB_UserFund(t *testing.T) {
	require := require.New(t)
	mockMap.UsersMap[tableUsers[1].requestUser.ID] = tableUsers[1].requestUser
	require.NoError(mockMap.UserFund(tableUsers[1].requestUser.ID, 300))
	require.Equal(tableUsers[1].resultUser.Balance, tableUsers[1].requestUser.Balance)
}

func TestDB_UserTake(t *testing.T) {
	require := require.New(t)
	mockMap.UsersMap[tableUsers[2].requestUser.ID] = tableUsers[2].requestUser
	mockMap.UsersMap[tableUsers[2].resultUser.ID] = tableUsers[2].resultUser
	require.Error(mockMap.UserTake(tableUsers[2].resultUser.ID, 200))
	require.NoError(mockMap.UserTake(tableUsers[1].requestUser.ID, 100))
	require.Equal(tableUsers[2].requestUser.Balance, tableUsers[2].resultUser.Balance)
}

func TestNew(t *testing.T) {
	require.NotEmpty(t, New())
}
