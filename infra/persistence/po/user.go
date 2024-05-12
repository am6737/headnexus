package po

type User struct {
	ID           string `bson:"_id,omitempty"`
	Name         string `bson:"name"`
	Email        string `bson:"email"`
	Token        string `bson:"token"`
	Password     string `bson:"password"`
	Status       uint   `bson:"status"`
	Verification string `bson:"verification"`
	CreatedAt    int64  `bson:"created_at"`
	UpdatedAt    int64  `bson:"updated_at"`
	DeletedAt    int64  `bson:"deleted_at"`
}
