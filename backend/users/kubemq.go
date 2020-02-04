package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/kubemq-io/kubemq-go"
	"log"
	"time"
)

type KubeMQ struct {
	client *kubemq.Client
	cfg    *Config
}

func NewKubeMQClient(cfg *Config) (*KubeMQ, error) {
	client, err := kubemq.NewClient(context.Background(),
		kubemq.WithAddress(cfg.KubeMQHost, cfg.KubeMQPort),
		kubemq.WithClientId(uuid.New().String()),
		kubemq.WithTransportType(kubemq.TransportTypeGRPC))
	if err != nil {
		return nil, err
	}
	k := &KubeMQ{
		client: client,
		cfg:    cfg,
	}
	return k, nil
}

func (k *KubeMQ) StartListenToCommands(ctx context.Context, channel, group string, commandsCh chan *kubemq.CommandReceive, errCh chan error) error {
	commandCh, err := k.client.SubscribeToCommands(ctx, channel, group, errCh)
	if err != nil {
		return err
	}
	go func() {
		for {
			select {
			case command, more := <-commandCh:
				if !more {
					return
				}
				commandsCh <- command

			case <-ctx.Done():
				return
			}

		}
	}()
	return nil
}

func (k *KubeMQ) StartListenToQueries(ctx context.Context, channel, group string, queryCh chan *kubemq.QueryReceive, errCh chan error) error {
	queriesCh, err := k.client.SubscribeToQueries(ctx, channel, group, errCh)
	if err != nil {
		return err
	}
	go func() {
		for {
			select {
			case query, more := <-queriesCh:
				if !more {
					return
				}
				queryCh <- query

			case <-ctx.Done():
				return
			}

		}
	}()
	return nil
}

func (k *KubeMQ) SendResponse(ctx context.Context, response *kubemq.Response) error {
	return k.client.R().
		SetBody(response.Body).
		SetMetadata(response.Metadata).
		SetError(response.Err).
		SetExecutedAt(response.ExecutedAt).
		SetResponseTo(response.ResponseTo).
		SetRequestId(response.RequestId).
		Send(ctx)

}

func (k *KubeMQ) SendToHistory(ctx context.Context, his *History) {

	err := k.client.E().SetChannel(k.cfg.HistoryChannel).SetId(his.Id).SetBody(his.Data()).SetMetadata(his.Method).Send(ctx)
	if err != nil {
		log.Println(fmt.Sprintf("error sending to history, error: %s", err.Error()))
	}
}

func (k *KubeMQ) SendToNotification(ctx context.Context, metadata string, data []byte) {
	res, err := k.client.ES().SetChannel(k.cfg.NotificationChannel).SetMetadata(metadata).SetBody([]byte(fmt.Sprintf("%s", string(data)))).Send(ctx)
	if err != nil {
		log.Println(fmt.Sprintf("error sending to notification, error: %s", err.Error()))
		return
	}
	log.Println(fmt.Sprintf("sending to notification, sent: %t, data: %s", res.Sent, data))

}

func (k *KubeMQ) SendCommandToCache(ctx context.Context, cm *CacheMessage) error {
	cr, err := k.client.C().SetChannel(k.cfg.CacheChannel).SetId(uuid.New().String()).SetBody(cm.Data()).SetTimeout(1 * time.Second).Send(ctx)
	if err != nil {
		return err
	}
	if !cr.Executed {
		return errors.New(cr.Error)
	}
	return nil
}
func (k *KubeMQ) SendQueryToCache(ctx context.Context, cm *CacheMessage) ([]byte, error) {
	qr, err := k.client.Q().SetChannel(k.cfg.CacheChannel).SetId(uuid.New().String()).SetBody(cm.Data()).SetTimeout(1 * time.Second).Send(ctx)
	if err != nil {
		return nil, err
	}
	if !qr.Executed {
		return nil, errors.New(qr.Error)
	}
	return qr.Body, nil
}
func (k *KubeMQ) Close() {
	_ = k.client.Close()
}
