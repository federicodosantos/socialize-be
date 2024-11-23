package query

const (
	InsertCommentQuery = `INSERT INTO comments(user_id, post_id, comment, created_at)
		VALUES(?, ?, ?, ?)`
	
	GetAllCommentsByPostIdQuery = `
	SELECT 
		c.id,
		c.post_id,
		c.user_id,
		c.comment,
		c.created_at,
		u.name AS user_name,      
		u.photo AS user_photo     
	FROM comments AS c
	JOIN users AS u ON u.id = c.user_id
	WHERE c.post_id = ?`

	DeleteCommentQuery = `DELETE FROM comments where id = ?`
)