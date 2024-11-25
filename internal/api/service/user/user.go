package user

import (
	"github.com/big-dust/DreamBridge/internal/api/types"
	"github.com/big-dust/DreamBridge/internal/model/user"
	"strings"
)

func SetUserInfo(id int, req *types.UserSetInfoReq) error {
	u := &user.User{
		Province:     req.Province,
		ExamType:     req.ExamType,
		SchoolType:   req.SchoolType,
		Physics:      req.Subject.Physics,
		History:      req.Subject.History,
		Chemistry:    req.Subject.Chemistry,
		Biology:      req.Subject.Biology,
		Geography:    req.Subject.Geography,
		Politics:     req.Subject.Politics,
		Score:        req.Score,
		ProvinceRank: req.ProvinceRank,
		Holland:      req.Holland,
		Interests:    strings.Join(req.Interests, " "),
	}
	return user.UpdateOne(id, u)
}
