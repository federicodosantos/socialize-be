package query

const (
	InsertUserQuery = `INSERT INTO users(id, name, email, password, created_at, updated_at) 
					VALUES (?, ?, ?, ?, ?, ?)`
	
	GetUserByIdQuery = `SELECT * FROM users WHERE id = ?`

	GetUserByEmailQuery = `SELECT * FROM users WHERE email = ?`

	UpdateUserDataQuery = `UPDATE users 
	SET name = ?, email = ?, password = ?, updated_at = ?
	WHERE id = ?`

	UpdateUserPhotoQuery = `UPDATE users 
	SET photo = ?, updated_at = ?
	WHERE id = ?`

	CheckEmailExistQuery = `SELECT COUNT(*) FROM users WHERE email = ?`
)