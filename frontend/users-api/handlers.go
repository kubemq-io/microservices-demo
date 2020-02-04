package main

import (
	"errors"
	"github.com/kubemq/demo/frontend/users-api/log"
	"github.com/labstack/echo"
)
var logger = log.NewLogger()

func (s *Server) register(c echo.Context) error {
	log:=logger.With("component","register")
	r := NewResponse(c, s.kube, "register", "query")
	acc := &NewUser{}
	err := c.Bind(acc)
	if err != nil {
		return r.SetError(err).Send()
	}
	r.SetRequestBody(acc)
	log.Infof("register: %s, %s, %s",acc.Name,acc.Password,acc.Email)
	resp, err := s.kube.SendQuery(c.Request().Context(), s.cfg.UsersChannel, "register", acc)
	if err != nil {
		log.Errorf("register account failed, error: %s",err.Error())
		return r.SetError(err).Send()
	}
	if !resp.Executed {
		log.Errorf("register account not executed, error: %s",resp.Error)
		return r.SetError(errors.New(resp.Error)).Send()
	}
	user, err := getUser(resp.Body)
	if err != nil {
		return r.SetError(err).Send()
	}

	r.SetResponseBody(user)
	return r.Send()
}

func (s *Server) login(c echo.Context) error {
	log:=logger.With("component","login")
	r := NewResponse(c, s.kube, "login", "query")
	lr := &LoginRequest{}
	err := c.Bind(lr)
	if err != nil {
		return r.SetError(err).Send()
	}
	log.Infof("login: %s, %s",lr.Name,lr.Password)
	r.SetRequestBody(lr)
	resp, err := s.kube.SendQuery(c.Request().Context(), s.cfg.UsersChannel, "login", lr)
	if err != nil {
		log.Errorf("login failed, error: %s",err.Error())
		return r.SetError(err).Send()
	}
	if !resp.Executed {
		log.Errorf("login not executed, error: %s",resp.Error)
		return r.SetError(errors.New(resp.Error)).Send()
	}
	loginResp, err := getLoginResponse(resp.Body)
	if err != nil {
		return r.SetError(err).Send()
	}

	r.SetResponseBody(loginResp)
	return r.Send()
}

func (s *Server) verify(c echo.Context) error {
	log:=logger.With("component","verify")
	r := NewResponse(c, s.kube, "verify", "command")
	vr := &VerifyRegistration{}
	err := c.Bind(vr)
	if err != nil {
		return r.SetError(err).Send()
	}
	r.SetRequestBody(vr)
	resp, err := s.kube.SendCommand(c.Request().Context(), s.cfg.UsersChannel, "verify_registration", vr)
	if err != nil {
		log.Errorf("verify failed, error: %s",err.Error())
		return r.SetError(err).Send()
	}
	if !resp.Executed {
		log.Errorf("verify not executed, Error: %s",resp.Error)
		return r.SetError(errors.New(resp.Error)).Send()
	}
	return r.Send()
}

func (s *Server) logout(c echo.Context) error {
	log:=logger.With("component","logout")
	r := NewResponse(c, s.kube, "logout", "command")
	lo := &LogoutRequest{}
	err := c.Bind(lo)
	if err != nil {
		return r.SetError(err).Send()
	}
	r.SetRequestBody(lo)
	resp, err := s.kube.SendCommand(c.Request().Context(), s.cfg.UsersChannel, "logout", lo)
	if err != nil {
		log.Errorf("logout failed, error: %s",err.Error())
		return r.SetError(err).Send()
	}
	if !resp.Executed {
		log.Errorf("logout not executed: %s",resp.Error)
		return r.SetError(errors.New(resp.Error)).Send()
	}
	return r.Send()
}

//
//func (s *Server) passwordReset(c echo.Context) error {
//	r := NewResponse(c)
//	pr := &PasswordResetRequest{}
//	err := c.Bind(pr)
//	if err != nil {
//		return r.SetError(err).Send()
//	}
//	r.SetRequestBody(pr)
//	resp, err := s.kube.SendCommand(c.Request().Context(), s.cfg.UsersChannel, "password_reset_request", pr)
//	if err != nil {
//		return r.SetError(err).Send()
//	}
//	if !resp.Executed {
//		return r.SetError(errors.New(resp.Error)).Send()
//	}
//	return r.Send()
//}

//
//func (s *Server) passwordChange(c echo.Context) error {
//	r := NewResponse(c)
//	pc := &PasswordChangeRequest{}
//	err := c.Bind(pc)
//	if err != nil {
//		return r.SetError(err).Send()
//	}
//	r.SetRequestBody(pc)
//	resp, err := s.kube.SendCommand(c.Request().Context(), s.cfg.UsersChannel, "password_change_request", pc)
//	if err != nil {
//		return r.SetError(err).Send()
//	}
//	if !resp.Executed {
//		return r.SetError(errors.New(resp.Error)).Send()
//	}
//	return r.Send()
//}
//
//func (s *Server) lock(c echo.Context) error {
//	r := NewResponse(c)
//	lc := &LockRequest{}
//	err := c.Bind(lc)
//	if err != nil {
//		return r.SetError(err).Send()
//	}
//	r.SetRequestBody(lc)
//	resp, err := s.kube.SendCommand(c.Request().Context(), s.cfg.UsersChannel, "lock", lc)
//	if err != nil {
//		return r.SetError(err).Send()
//	}
//	if !resp.Executed {
//		return r.SetError(errors.New(resp.Error)).Send()
//	}
//	return r.Send()
//}
//
//func (s *Server) unlock(c echo.Context) error {
//	r := NewResponse(c)
//	lc := &UnlockRequest{}
//	err := c.Bind(lc)
//	if err != nil {
//		return r.SetError(err).Send()
//	}
//	r.SetRequestBody(lc)
//	resp, err := s.kube.SendCommand(c.Request().Context(), s.cfg.UsersChannel, "unlock", lc)
//	if err != nil {
//		return r.SetError(err).Send()
//	}
//	if !resp.Executed {
//		return r.SetError(errors.New(resp.Error)).Send()
//	}
//	return r.Send()
//}
