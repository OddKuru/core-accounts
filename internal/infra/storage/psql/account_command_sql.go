package psql

const (
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
)
