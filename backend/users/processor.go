package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/kubemq-io/kubemq-go"
	"log"
	"time"
)

type Processor struct {
	pq    *Postgres
	kube  *KubeMQ
	cfg   *Config
	cache *Cache
}

func NewProcessor(ctx context.Context, pq *Postgres, kube *KubeMQ, cfg *Config) *Processor {
	p := &Processor{
		pq:    pq,
		kube:  kube,
		cfg:   cfg,
		cache: NewCahce(kube),
	}
	go p.run(ctx)
	return p
}

func (p Processor) run(ctx context.Context) {
	commandsCh := make(chan *kubemq.CommandReceive, 1)
	queriesCh := make(chan *kubemq.QueryReceive, 1)
	errCh := make(chan error, 1)

	for {
		err := p.kube.StartListenToCommands(ctx, p.cfg.UsersChannel, p.cfg.Group, commandsCh, errCh)
		if err != nil {
			log.Printf("error connecting to kubemq, error: %s, retrying...\n", err.Error())
			time.Sleep(time.Second)
		} else {
			break
		}
	}
	for {
		err := p.kube.StartListenToQueries(ctx, p.cfg.UsersChannel, p.cfg.Group, queriesCh, errCh)
		if err != nil {
			log.Printf("error connecting to kubemq, error: %s, retrying...\n", err.Error())
			time.Sleep(time.Second)
		} else {
			break
		}
	}
	log.Println("waiting for commands / queries from KubeMQ")

	for {
		select {
		case command := <-commandsCh:
			var err error
			switch command.Metadata {
			case "verify_registration":
				err = p.verifyRegistration(ctx, command.Body)
			case "logout":
				err = p.logout(ctx, command.Body)

			case "password_reset_request":
				err = p.passwordResetRequest(ctx, command.Body)
			case "password_change_request":
				err = p.passwordChangeRequest(ctx, command.Body)
			case "lock":
				err = p.lock(ctx, command.Body)
			case "unlock":
				err = p.unlock(ctx, command.Body)
			default:
				err = errors.New("unknown command")
			}

			resp := &kubemq.Response{
				RequestId:  command.Id,
				ResponseTo: command.ResponseTo,
				Metadata:   command.Metadata,
				Body:       nil,
			}
			if err != nil {
				resp.Err = err
			} else {
				resp.ExecutedAt = time.Now()
			}
			err = p.kube.SendResponse(ctx, resp)
			his := &History{
				Id:           command.Id,
				Source:       "users-service",
				Time:         time.Now(),
				Type:         "command",
				Method:       command.Metadata,
				Request:      fmt.Sprintf("%s", command.Body),
				Response:     "",
				IsError:      false,
				ErrorMessage: "",
			}
			if err != nil {
				his.IsError = true
				his.ErrorMessage = err.Error()
			}
			go p.kube.SendToHistory(ctx, his)
			log.Println(fmt.Sprintf("command: %s, data: %s, error: %s", command.Metadata, command.Body, resp.Err))
		case query := <-queriesCh:
			var resp []byte
			var err error
			switch query.Metadata {
			case "register":
				resp, err = p.register(ctx, query.Body)
			case "login":
				resp, err = p.login(ctx, query.Body)

			}
			response := &kubemq.Response{
				RequestId:  query.Id,
				ResponseTo: query.ResponseTo,
				Metadata:   query.Metadata,
			}
			if err != nil {
				response.Err = err
			} else {
				response.ExecutedAt = time.Now()
				response.Body = resp
			}
			err = p.kube.SendResponse(ctx, response)
			his := &History{
				Id:           query.Id,
				Source:       "users-service",
				Time:         time.Now(),
				Type:         "query",
				Method:       query.Metadata,
				Request:      fmt.Sprintf("%s", query.Body),
				Response:     "",
				IsError:      false,
				ErrorMessage: "",
			}
			if response.Err != nil {
				his.IsError = true
				his.ErrorMessage = response.Err.Error()
			} else {
				his.Response = fmt.Sprintf("%s", response.Body)
			}
			go p.kube.SendToHistory(ctx, his)
			log.Println(fmt.Sprintf("query: %s, data: %s, response: %s,error: %s", query.Metadata, query.Body, response.Body, response.Err))

		case err := <-errCh:
			log.Fatal(err)
		case <-ctx.Done():
			return
		}

	}

}

func (p *Processor) register(ctx context.Context, data []byte) ([]byte, error) {
	newUser, err := getNewUser(data)
	if err != nil {
		return nil, err
	}
	user, err := p.pq.Register(ctx, newUser)
	if err != nil {
		return nil, err
	}
	go p.kube.SendToNotification(ctx, "register", []byte(PrettyJson(user)))
	return user.Data(), nil
}

func (p *Processor) verifyRegistration(ctx context.Context, data []byte) error {
	vr, err := getVerificationRegistration(data)
	if err != nil {
		return err
	}
	err = p.pq.VerifyRegistration(ctx, vr)
	if err != nil {
		return err
	}

	return nil
}

func (p *Processor) login(ctx context.Context, data []byte) ([]byte, error) {
	lr, err := getLoginRequest(data)
	if err != nil {
		return nil, err
	}
	respDataFromCache, err := p.cache.Get(ctx, lr.Name)
	if err == nil {
		lr := &LoginResponse{}
		err := json.Unmarshal(respDataFromCache, lr)
		fmt.Println(err, lr)
		return respDataFromCache, nil
	}
	resp, err := p.pq.Login(ctx, lr)
	if err != nil {
		return nil, err
	}

	go p.cache.Set(ctx, lr.Name, resp, resp.Expiry)
	return resp.Data(), nil
}

func (p *Processor) logout(ctx context.Context, data []byte) error {
	lr, err := getLogoutRequest(data)
	if err != nil {
		return err
	}
	logoutResp, err := p.pq.Logout(ctx, lr)
	if err != nil {
		return err
	}
	go p.cache.Del(ctx, logoutResp.Name)
	return nil
}

func (p *Processor) passwordResetRequest(ctx context.Context, data []byte) error {
	prr, err := getPasswordResetRequest(data)
	if err != nil {
		return err
	}
	pr, err := p.pq.PasswordResetRequest(ctx, prr)
	if err != nil {
		return err
	}

	go p.kube.SendToNotification(ctx, "password_reset_request", []byte(PrettyJson(pr)))
	return nil
}

func (p *Processor) passwordChangeRequest(ctx context.Context, data []byte) error {
	prr, err := getPasswordChaneRequest(data)
	if err != nil {
		return err
	}
	pr, err := p.pq.PasswordChangeRequest(ctx, prr)
	if err != nil {
		return err
	}

	go p.kube.SendToNotification(ctx, "password_change_request", []byte(PrettyJson(pr)))
	return nil
}

func (p *Processor) lock(ctx context.Context, data []byte) error {
	lr, err := getLockRequest(data)
	if err != nil {
		return err
	}
	lres, err := p.pq.Lock(ctx, lr)
	if err != nil {
		return err
	}

	go p.kube.SendToNotification(ctx, "lock", []byte(PrettyJson(lres)))

	return nil
}

func (p *Processor) unlock(ctx context.Context, data []byte) error {
	ur, err := getUnlockRequest(data)
	if err != nil {
		return err
	}
	lres, err := p.pq.Unlock(ctx, ur)
	if err != nil {
		return err
	}
	go p.kube.SendToNotification(ctx, "unlock", []byte(PrettyJson(lres)))
	return nil
}
