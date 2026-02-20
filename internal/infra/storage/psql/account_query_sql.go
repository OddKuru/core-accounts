package psql

const (
	AccountHasByEmailQuery = `SELECT exists(SELECT 1 FROM accounts WHERE email = $1);`
	AccountHasByNameQuery  = `SELECT exists(SELECT 1 FROM accounts WHERE login = $1);`
	AccountGetByIdQuery    = `
		SELECT 
		    a.login,
		    a.email,
		    a.password_hash, 
		    a.version,
		    a.created_at,
		    a.updated_at, 
		    r.value AS role 
		FROM accounts a
		JOIN roles r ON r.id = a.role_id
		WHERE a.id = $1;
	`
)
