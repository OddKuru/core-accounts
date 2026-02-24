package psql

const (
	AccountHasByIDQuery = `
		SELECT exists(SELECT 1 FROM accounts WHERE id = $1);
	`
	AccountCreateQuery = `
		INSERT INTO accounts (id, login, email, password_hash, version, role_id, updated_at, created_at) 
		SELECT 
			$1,
			$2,
			$3,
			$4,
			$5,
			r.id,
			$7,
			$8
		FROM roles r
		WHERE r.value = $6;
	`
	AccountUpdateQuery = `
		UPDATE accounts AS a
		SET
			login = $2,
			email = $3,
			password_hash = $4,
			version = a.version+1,
			role_id = r.id,
			updated_at = $5
		FROM roles AS r
		WHERE r.value = $6
			AND a.id = $1
			AND a.version = $7;
	`
)
