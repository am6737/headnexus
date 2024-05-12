package user

import (
	"context"
	"github.com/am6737/headnexus/domain/user/entity"
	"github.com/am6737/headnexus/pkg/email"
	pkgjwt "github.com/am6737/headnexus/pkg/jwt"
)

type Application struct {
	Commands    Commands
	Queries     Queries
	JwtConfig   *pkgjwt.JWTConfig
	EmailClient *email.EmailClient
}

type Commands struct {
	Handler CommandHandler
}

type Queries struct {
	Handler QueryHandler
}

type CreateUser struct {
	Name         string `json:"name"`
	Email        string `json:"email"`
	Password     string `json:"password"`
	Verification string `json:"verification"`
	Status       uint   `json:"status"`
	Token        string `json:"token"`
}

type DeleteUser struct {
	ID string
}

type UpdateUser struct {
	ID           string
	Name         string `json:"name"`
	Email        string `json:"email"`
	Password     string `json:"password"`
	Verification string `json:"verification"`
	Status       uint   `json:"status"`
	Token        string `json:"token"`
}

type CommandHandler interface {
	Create(ctx context.Context, cmd *CreateUser) (*entity.User, error)
	Delete(ctx context.Context, cmd *DeleteUser) error
	Update(ctx context.Context, cmd *UpdateUser) (*entity.User, error)
	Register(ctx context.Context, cmd *CreateUser) (*entity.User, error)
	Login(ctx context.Context, email string, password string) (string, error)
}

type QueryHandler interface {
	Get(ctx context.Context, query *GetUser) (*entity.User, error)
	Find(ctx context.Context, query *FindUser) ([]*entity.User, error)
}

type GetUser struct {
	ID string
}

type FindUser struct {
	Name         string
	Email        string
	Token        string
	Verification string
}

//
//type FindRule struct {
//	Name     string
//	HostID   string
//	PageSize int
//	PageNum  int
//}
//
//func (fr FindRule) String() string {
//	marshal, err := json.Marshal(fr)
//	if err != nil {
//		return ""
//	}
//	return string(marshal)
//}
//
