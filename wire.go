// +build wireinject

package main

import (
	"context"
	"go-chat/app/pkg/filesystem"
	"go-chat/app/repository"
	"go-chat/provider"

	"github.com/google/wire"
	"go-chat/app/cache"
	"go-chat/app/http/handler"
	"go-chat/app/http/handler/api/v1"
	"go-chat/app/http/handler/open"
	"go-chat/app/http/handler/ws"
	"go-chat/app/http/router"
	"go-chat/app/service"
	"go-chat/config"
)

var providerSet = wire.NewSet(
	// 基础服务
	provider.NewLogger,
	provider.RedisConnect,
	provider.MysqlConnect,
	provider.NewHttp,

	// 注册路由
	router.NewRouter,

	// other
	filesystem.NewFilesystem,

	// 缓存
	cache.NewServerRun,
	wire.Struct(new(cache.WsClient), "*"),
	wire.Struct(new(cache.AuthTokenCache), "*"),
	wire.Struct(new(cache.SmsCodeCache), "*"),
	wire.Struct(new(cache.RedisLock), "*"),

	// repo
	wire.Struct(new(repository.UserRepository), "*"),
	wire.Struct(new(repository.TalkRecordsRepo), "*"),
	wire.Struct(new(repository.TalkRecordsCodeRepo), "*"),
	wire.Struct(new(repository.TalkRecordsLoginRepo), "*"),
	wire.Struct(new(repository.TalkRecordsFileRepo), "*"),
	wire.Struct(new(repository.TalkRecordsVoteRepo), "*"),

	// 服务
	service.NewUserService,
	service.NewSmsService,
	service.NewTalkMessageService,
	service.NewClientService,
	service.NewGroupService,
	service.NewGroupMemberService,
	service.NewGroupNoticeService,
	service.NewTalkListService,
	wire.Struct(new(service.SocketService), "*"),

	// handler 处理
	v1.NewAuthHandler,
	v1.NewCommonHandler,
	v1.NewUserHandler,
	v1.NewGroupHandler,
	v1.NewGroupNoticeHandler,
	v1.NewTalkHandler,
	v1.NewTalkMessageHandler,
	v1.NewUploadHandler,
	v1.NewDownloadHandler,
	v1.NewEmoticonHandler,
	open.NewIndexHandler,
	ws.NewWebSocketHandler,

	wire.Struct(new(handler.Handler), "*"),
	wire.Struct(new(provider.Services), "*"),
)

func Initialize(ctx context.Context, conf *config.Config) *provider.Services {
	panic(wire.Build(providerSet))
}
