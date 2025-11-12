package domain

type User struct {
    ID          string `bson:"_id,omitempty" json:"id,omitempty"`
    Name        string `bson:"name" json:"name"`
    Email       string `bson:"email" json:"email"`
    PhoneNumber string `bson:"phoneNumber" json:"phoneNumber"`
    Password    string `bson:"password" json:"-"`
    ProfileImage string `bson:"profileImage,omitempty" json:"profileImage,omitempty"`
    CreatedAt   int64  `bson:"createdAt" json:"createdAt"`
    UpdatedAt   int64  `bson:"updatedAt" json:"updatedAt"`
}

type UserRepository interface {
    Create(user *User) error
    FindByEmail(email string) (*User, error)
    FindByID(id string) (*User, error)
    Update(user *User) error
    Delete(id string) error
}

type UserService interface {
    Register(user *User) error
    Login(email, password string) (string, error)
    GetProfile(id string) (*User, error)
    UpdateProfile(user *User) error
    ChangePassword(userID, currentPassword, newPassword string) error
}
