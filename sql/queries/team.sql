-- name: NewTeam :one
INSERT INTO teams (
		name,
		color
) VALUES (
		@name,
    @color
) RETURNING *;

-- name: AddTeamPermissions :exec
INSERT INTO team_permissions ( team_id, user_id,is_owner) VALUES ( @teamId, @userId, @isOwner);

-- name: InviteToTeam :exec
INSERT INTO team_invitations (team_id, sender_id, receiver_id) VALUES(@teamId,@senderId,@receiverId);

-- name: ResondToTeamInvitation :exec 
UPDATE team_invitations SET status = @status where id = @invitationId;
