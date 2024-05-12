package converter

import (
	"github.com/am6737/headnexus/domain/user/entity"
	"github.com/am6737/headnexus/infra/persistence/po"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func UserEntityToPO(e *entity.User) (*po.User, error) {
	m := &po.User{}
	if e.ID != "" {
		oid, err := primitive.ObjectIDFromHex(e.ID)
		if err != nil {
			return nil, err
		}
		m.ID = oid.Hex()
	}
	m.Name = e.Name
	m.Email = e.Email
	m.Token = e.Token
	m.Status = uint(e.Status)
	m.Verification = e.Verification
	m.Password = e.Password
	return m, nil
}

func UserPOToEntity(po *po.User) (*entity.User, error) {
	return &entity.User{
		ID:           po.ID,
		Name:         po.Name,
		Email:        po.Email,
		Token:        po.Token,
		Password:     po.Password,
		Status:       entity.UserStatus(po.Status),
		Verification: po.Verification,
	}, nil
}
