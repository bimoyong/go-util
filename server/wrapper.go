package server

import (
	"context"
	"fmt"

	"github.com/micro/go-micro/v2/client"
	log "github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2/metadata"
	"github.com/micro/go-micro/v2/server"

	umetadata "github.com/bimoyong/go-util/metadata"
	"gitlab.com/bimoyong/go-user/auth"
	puser "gitlab.com/bimoyong/go-user/proto/user"
)

// LogWrapper function
func LogWrapper(fn server.SubscriberFunc) server.SubscriberFunc {
	return func(ctx context.Context, msg server.Message) (err error) {
		log.Debugf("Consume message: topic=[%s] header=[%+v] data=[%s]",
			msg.Topic(), msg.Header(), msg.Body())

		return fn(ctx, msg)
	}
}

// AuthWrapper function
func AuthWrapper(fn server.SubscriberFunc) server.SubscriberFunc {
	return func(ctx context.Context, msg server.Message) (err error) {
		srv := puser.NewUserService("go.srv.user", client.DefaultClient)
		var req puser.EmptyReq
		var rsp *puser.InspectRsp
		// TODO: inspect authorization with timeout
		md, _ := metadata.FromContext(ctx)
		ctx = metadata.NewContext(context.Background(), umetadata.NewMetadata(md))
		tk, _ := md.Get("Authorization")
		if rsp, err = srv.Inspect(ctx, &req); err != nil {
			log.Warnf("Bad credential: authz=[%s]", tk)

			return
		}
		log.Debugf("Login success: authz=[%s] user=[%s]", tk, rsp.Id)

		var kind auth.Kind
		if rsp.Id == "system" {
			kind = auth.System
		} else {
			kind = auth.Registered
		}

		ctx = metadata.Set(ctx, "User", rsp.Id)
		ctx = metadata.Set(ctx, "Kind", fmt.Sprintf("%d", kind))

		return fn(ctx, msg)
	}
}
