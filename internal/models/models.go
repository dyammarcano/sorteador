package models

type Sponsor struct {
	UUID      string `db:"uuid" json:"uuid"`
	Name      string `db:"name" json:"name"`
	Prize     string `db:"prize" json:"prize"`
	Timestamp int64  `db:"timestamp" json:"timestamp"`
	ULID      string `db:"ulid" json:"ulid"`
}

type Participant struct {
	Name  string `db:"name" json:"name"`
	Email string `db:"email" json:"email"`
}
