// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//+build !wireinject

package main

import (
	"context"
	"github.com/google/wire"
	"go-chat/app/cache"
	"go-chat/app/dao"
	"go-chat/app/http/handler"
	"go-chat/app/http/handler/api/v1"
	"go-chat/app/http/handler/open"
	"go-chat/app/http/handler/ws"
	"go-chat/app/http/router"
	"go-chat/app/pkg/filesystem"
	"go-chat/app/process"
	"go-chat/app/service"
	"go-chat/provider"
)

import (
	_ "go-chat/app/pkg/validation"
)

// Injectors from wire.go:

func Initialize(ctx context.Context) *provider.Services {
	config := provider.NewConfig()
	client := provider.NewRedisClient(ctx, config)
	smsCodeCache := &cache.SmsCodeCache{
		Redis: client,
	}
	smsService := service.NewSmsService(smsCodeCache)
	db := provider.NewMySQLClient(config)
	base := dao.NewBaseDao(db)
	userDao := &dao.UserDao{
		Base: base,
	}
	common := v1.NewCommonHandler(config, smsService, userDao)
	userService := service.NewUserService(userDao)
	session := cache.NewSession(client)
	redisLock := cache.NewRedisLock(client)
	auth := v1.NewAuthHandler(config, userService, smsService, session, redisLock)
	user := v1.NewUserHandler(userService, smsService)
	baseService := service.NewBaseService(db, client)
	groupMemberService := service.NewGroupMemberService(db)
	unreadTalkCache := cache.NewUnreadTalkCache(client)
	talkMessageForwardService := service.NewTalkMessageForwardService(baseService)
	lastMessage := cache.NewLastMessage(client)
	talkMessageService := service.NewTalkMessageService(baseService, config, groupMemberService, unreadTalkCache, talkMessageForwardService, lastMessage)
	talkService := service.NewTalkService(baseService, groupMemberService)
	talkMessage := v1.NewTalkMessageHandler(talkMessageService, talkService)
	talkListDao := dao.NewTalkListDao(base)
	talkListService := service.NewTalkListService(baseService, talkListDao)
	serverRunID := cache.NewServerRun(client)
	wsClientSession := cache.NewWsClientSession(client, config, serverRunID)
	usersFriendsDao := dao.NewUsersFriends(base, client)
	talk := v1.NewTalkHandler(talkService, talkListService, redisLock, userService, wsClientSession, lastMessage, usersFriendsDao)
	talkRecordsService := service.NewTalkRecordsService(baseService)
	talkRecords := v1.NewTalkRecordsHandler(talkRecordsService)
	download := v1.NewDownloadHandler()
	emoticonDao := dao.NewEmoticonDao(base)
	emoticonService := service.NewEmoticonService(baseService, emoticonDao)
	filesystemFilesystem := filesystem.NewFilesystem(config)
	emoticon := v1.NewEmoticonHandler(emoticonService, filesystemFilesystem, redisLock)
	upload := v1.NewUploadHandler(config, filesystemFilesystem)
	index := open.NewIndexHandler(client)
	clientService := service.NewClientService(wsClientSession)
	defaultWebSocket := ws.NewDefaultWebSocket(clientService)
	groupDao := &dao.GroupDao{
		Base: base,
	}
	groupService := service.NewGroupService(groupDao, db, groupMemberService)
	group := v1.NewGroupHandler(groupService, groupMemberService, talkListService, userDao, redisLock)
	groupNoticeDao := &dao.GroupNoticeDao{
		Base: base,
	}
	groupNoticeService := service.NewGroupNoticeService(groupNoticeDao)
	groupNotice := v1.NewGroupNoticeHandler(groupNoticeService, groupMemberService)
	handlerHandler := &handler.Handler{
		Common:           common,
		Auth:             auth,
		User:             user,
		TalkMessage:      talkMessage,
		Talk:             talk,
		TalkRecords:      talkRecords,
		Download:         download,
		Emoticon:         emoticon,
		Upload:           upload,
		Index:            index,
		DefaultWebSocket: defaultWebSocket,
		Group:            group,
		GroupNotice:      groupNotice,
	}
	engine := router.NewRouter(config, handlerHandler, session)
	server := provider.NewHttpServer(config, engine)
	serverRun := process.NewServerRun(config, serverRunID)
	wsSubscribe := process.NewWsSubscribe(client)
	processProcess := process.NewProcessManage(serverRun, wsSubscribe)
	services := &provider.Services{
		Config:     config,
		HttpServer: server,
		Process:    processProcess,
	}
	return services
}

// wire.go:

var providerSet = wire.NewSet(provider.NewConfig, provider.NewLogger, provider.NewMySQLClient, provider.NewRedisClient, provider.NewHttpClient, provider.NewHttpServer, router.NewRouter, filesystem.NewFilesystem, cache.NewSession, cache.NewServerRun, cache.NewUnreadTalkCache, cache.NewRedisLock, cache.NewWsClientSession, cache.NewLastMessage, wire.Struct(new(cache.SmsCodeCache), "*"), dao.NewBaseDao, dao.NewUsersFriends, wire.Struct(new(dao.UserDao), "*"), wire.Struct(new(dao.TalkRecordsDao), "*"), wire.Struct(new(dao.TalkRecordsCodeDao), "*"), wire.Struct(new(dao.TalkRecordsLoginDao), "*"), wire.Struct(new(dao.TalkRecordsFileDao), "*"), wire.Struct(new(dao.TalkRecordsVoteDao), "*"), wire.Struct(new(dao.GroupDao), "*"), wire.Struct(new(dao.GroupNoticeDao), "*"), dao.NewTalkListDao, dao.NewEmoticonDao, service.NewBaseService, service.NewUserService, service.NewSmsService, service.NewTalkService, service.NewTalkMessageService, service.NewClientService, service.NewGroupService, service.NewGroupMemberService, service.NewGroupNoticeService, service.NewTalkListService, service.NewTalkMessageForwardService, service.NewEmoticonService, service.NewTalkRecordsService, v1.NewAuthHandler, v1.NewCommonHandler, v1.NewUserHandler, v1.NewGroupHandler, v1.NewGroupNoticeHandler, v1.NewTalkHandler, v1.NewTalkMessageHandler, v1.NewUploadHandler, v1.NewDownloadHandler, v1.NewEmoticonHandler, v1.NewTalkRecordsHandler, open.NewIndexHandler, ws.NewDefaultWebSocket, wire.Struct(new(handler.Handler), "*"), wire.Struct(new(provider.Services), "*"), process.NewWsSubscribe, process.NewServerRun, process.NewProcessManage)
