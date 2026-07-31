package projectmembers

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestServiceAddMemberValidatesAndStoresMembership(t *testing.T) {
	service := NewService()

	membership, err := service.AddMember(context.Background(), " project-1 ", " user-1 ", " admin ")

	require.NoError(t, err)
	assert.Equal(t, Membership{ProjectID: "project-1", UserID: "user-1", Role: "admin"}, membership)
}

func TestServiceListMembersReturnsSortedMembers(t *testing.T) {
	service := NewService()
	_, err := service.AddMember(context.Background(), "project-1", "user-2", "viewer")
	require.NoError(t, err)
	_, err = service.AddMember(context.Background(), "project-1", "user-1", "editor")
	require.NoError(t, err)

	members, err := service.ListMembers(context.Background(), "project-1")

	require.NoError(t, err)
	assert.Equal(t, []string{"user-1", "user-2"}, []string{members[0].UserID, members[1].UserID})
}

func TestServiceRemoveMemberDeletesExistingMembership(t *testing.T) {
	service := NewService()
	_, err := service.AddMember(context.Background(), "project-1", "user-1", "viewer")
	require.NoError(t, err)

	require.NoError(t, service.RemoveMember(context.Background(), "project-1", "user-1"))
	members, err := service.ListMembers(context.Background(), "project-1")
	require.NoError(t, err)
	assert.Empty(t, members)
}

func TestServiceRejectsInvalidMembershipInput(t *testing.T) {
	service := NewService()

	_, err := service.AddMember(context.Background(), "", "user-1", "viewer")
	assert.ErrorIs(t, err, ErrInvalidMembership)

	err = service.RemoveMember(context.Background(), "project-1", "")
	assert.ErrorIs(t, err, ErrInvalidMembership)

	_, err = service.ListMembers(context.Background(), "")
	assert.ErrorIs(t, err, ErrInvalidMembership)
}

func TestServiceHonorsCanceledContext(t *testing.T) {
	service := NewService()
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := service.AddMember(ctx, "project-1", "user-1", "viewer")

	assert.ErrorIs(t, err, context.Canceled)
}
