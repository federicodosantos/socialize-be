package query

const (
	InsertPostQuery = `
	INSERT INTO posts (
		user_id, title, content, image, created_at, updated_at
	) VALUES (
		?, ?, ?, ?, ?, ?
	)`

	GetAllPostQuery = `
	SELECT 
		p.id,
		p.title,
		p.content,
		p.user_id,
		p.image,
		u.name AS user_name,
		u.photo AS user_photo,
		p.created_at,
		p.updated_at, 
		(SELECT count(*) from votes WHERE vote = 1 AND post_id = p.id) AS up_vote, 
		(SELECT count(*) from votes WHERE vote = -1 AND post_id = p.id) AS down_vote  
	FROM posts AS p
	JOIN users AS u ON u.id = p.user_id
	`
	GetAllPostByUserIDQuery = `
	SELECT 
		p.id,
		p.title,
		p.content,
		p.user_id,
		p.image,
		u.name AS user_name,
		u.photo AS user_photo,
		p.created_at,
		p.updated_at, 
		(SELECT count(*) from votes WHERE vote = 1 AND post_id = p.id) AS up_vote, 
		(SELECT count(*) from votes WHERE vote = -1 AND post_id = p.id) AS down_vote   
	FROM posts AS p
	JOIN users AS u ON u.id = p.user_id
	WHERE p.user_id = ?`

	GetPostByIDQuery = `
	SELECT 
		p.id,
		p.title,
		p.content,
		p.user_id,
		p.image,
		u.name AS user_name,
		u.photo AS user_photo,
		p.created_at,
		p.updated_at, 
		(SELECT count(*) from votes WHERE vote = 1 AND post_id = p.id) AS up_vote, 
		(SELECT count(*) from votes WHERE vote = -1 AND post_id = p.id) AS down_vote   
	FROM posts AS p
	JOIN users AS u ON u.id = p.user_id 
	WHERE p.id = ?`

	DeletePostQuery = `DELETE FROM posts WHERE id = ?`

	CheckPostExistQuery = `SELECT COUNT(*) FROM posts WHERE id = ?`

	CheckVoteExistQuery = `SELECT COUNT(*) FROM votes WHERE post_id = ? AND user_id = ?`

	InsertVoteQuery = `
	INSERT INTO votes (post_id, user_id, vote)
	VALUES (?, ?, ?)
	`

	UpdateVoteQuery = `UPDATE votes SET vote = ? WHERE post_id = ? AND user_id = ?`

	DeleteVoteQuery = `DELETE FROM votes WHERE post_id = ? AND user_id = ?`
)