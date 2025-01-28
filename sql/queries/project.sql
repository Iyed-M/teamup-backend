-- name: NewProject :one
INSERT INTO projects (
		name,
		color
) VALUES (
		@Name,
    @Color
) RETURNING *;

-- name: AddUserToProject :exec
INSERT INTO user_projects ( project_id, user_id,is_owner) VALUES ( @ProjectId, @UserId, @IsOwner);

-- name: InviteToProject :exec
INSERT INTO project_invitations (project_id, sender_id, receiver_id) VALUES(@ProjectId,@SenderId,@ReceiverId);

-- name: ResondToProjectInvitation :exec 
UPDATE project_invitations SET status = @Status where id = @InvitationId;
