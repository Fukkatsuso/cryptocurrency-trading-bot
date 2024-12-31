package persistence

// type sessionRepository struct {
// 	db DB
// }

// func NewSessionRepository(db DB) repository.SessionRepository {
// 	return &sessionRepository{
// 		db: db,
// 	}
// }

// func (sr *sessionRepository) Save(userID string, sessionID string) error {
// 	cmd := `
//         UPDATE
//             users
//         SET
//             session_id_hash = ?
//         WHERE
//             id = ?
//     `
// 	_, err := sr.db.Exec(cmd, sessionID, userID)
// 	return err
// }

// func (sr *sessionRepository) FindByUserID(userID string) (string, error) {
// 	cmd := `
//         SELECT
//             session_id_hash
//         FROM
//             users
//         WHERE
//             id = ?
//     `
// 	row := sr.db.QueryRow(cmd, userID)

// 	var sessionID string
// 	err := row.Scan(&sessionID)
// 	if err == sql.ErrNoRows {
// 		return "", errors.New("user is not found")
// 	}
// 	if err != nil {
// 		return "", err
// 	}

// 	return sessionID, nil
// }

// func (sr *sessionRepository) Delete(userID string) error {
// 	cmd := `
//         UPDATE
//             users
//         SET
//             session_id_hash = ''
//         WHERE
//             id = ?
//     `
// 	_, err := sr.db.Exec(cmd, userID)
// 	return err
// }
