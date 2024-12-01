package usersUsecases

import (
	"fmt"
	"go_learn_project_rest_api/config"
	"go_learn_project_rest_api/modules/users"
	"go_learn_project_rest_api/modules/users/usersRepositories"
	"go_learn_project_rest_api/pkgs/auth"

	"golang.org/x/crypto/bcrypt"
)

type IUsersUsecases interface {
	InsertCustomer(*users.UserRegisterReq) (*users.UserPassport, error)
	GetPassport(*users.UserCredential) (*users.UserPassport, error)
	RefreshPassport(*users.UserRefreshCredential) (*users.UserPassport, error)
	DeleteOauth(string) error
	InsertAdmin(*users.UserRegisterReq) (*users.UserPassport, error)
	GetUserProfile(string) (*users.User, error)
}

type usersUsecases struct {
	cfg             config.IConfig
	usersRepository usersRepositories.IUsersRepository
}

func UsersUsecases(cfg config.IConfig, usersRepository usersRepositories.IUsersRepository) IUsersUsecases {
	return &usersUsecases{
		usersRepository: usersRepository,
		cfg:             cfg,
	}
}

func (u *usersUsecases) InsertCustomer(request *users.UserRegisterReq) (*users.UserPassport, error) {
	if err := request.BcryptHashing(); err != nil {
		return nil, err
	}

	result, err := u.usersRepository.InsertUser(request, false)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (u *usersUsecases) InsertAdmin(request *users.UserRegisterReq) (*users.UserPassport, error) {
	if err := request.BcryptHashing(); err != nil {
		return nil, err
	}

	result, err := u.usersRepository.InsertUser(request, true)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (u *usersUsecases) GetPassport(request *users.UserCredential) (*users.UserPassport, error) {
	user, err := u.usersRepository.FindOneUserByEmail(request.Email)
	if err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password)); err != nil {
		return nil, fmt.Errorf("password is invalid")
	}

	newToken, err := auth.NewAuth(auth.Access, u.cfg.Jwt(), &users.UserClaims{
		Id:     user.Id,
		RoleId: user.RoleId,
	})
	if err != nil {
		return nil, err
	}
	refreshToken, err := auth.NewAuth(auth.Refresh, u.cfg.Jwt(), &users.UserClaims{
		Id:     user.Id,
		RoleId: user.RoleId,
	})
	if err != nil {
		return nil, err
	}

	passport := &users.UserPassport{
		User: &users.User{
			Id:       user.Id,
			Email:    user.Email,
			RoleId:   user.RoleId,
			Username: user.Username,
		},
		Token: &users.UserToken{
			AccessToken:  newToken.SignToken(),
			RefreshToken: refreshToken.SignToken(),
		},
	}

	if err := u.usersRepository.InsertOauth(passport); err != nil {
		return nil, err
	}

	return passport, nil
}

func (u *usersUsecases) RefreshPassport(req *users.UserRefreshCredential) (*users.UserPassport, error) {
	claims, err := auth.ParseToken(u.cfg.Jwt(), req.RefreshToken)
	if err != nil {
		return nil, err
	}

	oauth, err := u.usersRepository.FindOneOauth(req.RefreshToken)
	if err != nil {
		return nil, err
	}

	profile, err := u.usersRepository.GetProfile(oauth.UserId)
	if err != nil {
		return nil, err
	}

	newClaims := &users.UserClaims{
		Id:     profile.Id,
		RoleId: profile.RoleId,
	}

	accessToken, err := auth.NewAuth(auth.Access, u.cfg.Jwt(), newClaims)
	if err != nil {
		return nil, err
	}

	refreshToken := auth.RepeatToken(u.cfg.Jwt(), newClaims, claims.ExpiresAt.Unix())

	passport := &users.UserPassport{
		User: profile,
		Token: &users.UserToken{
			Id:           oauth.Id,
			AccessToken:  accessToken.SignToken(),
			RefreshToken: refreshToken,
		},
	}

	if err := u.usersRepository.UpdateOauth(passport.Token); err != nil {
		return nil, err
	}

	return passport, nil
}

func (u *usersUsecases) DeleteOauth(oauthId string) error {
	if err := u.usersRepository.DeleteOauth(oauthId); err != nil {
		return err
	}
	return nil
}

func (u *usersUsecases) GetUserProfile(userId string) (*users.User, error) {
	return u.usersRepository.GetProfile(userId)
}
