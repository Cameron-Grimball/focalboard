package main

import (
	"fmt"

	"github.com/mattermost/focalboard/server/services/notify/notifymentions"
	"github.com/mattermost/focalboard/server/services/notify/notifysubscriptions"
	"github.com/mattermost/focalboard/server/services/notify/plugindelivery"
	"github.com/mattermost/focalboard/server/services/store"

	pluginapi "github.com/mattermost/mattermost-plugin-api"

	"github.com/mattermost/mattermost-server/v6/model"

	"github.com/mattermost/mattermost-server/v6/shared/mlog"
)

const (
	botUsername    = "boards"
	botDisplayname = "Boards"
	botDescription = "Created by Boards plugin."
)

type notifyBackendParams struct {
	client     *pluginapi.Client
	serverRoot string
	logger     *mlog.Logger
}

func createMentionsNotifyBackend(params notifyBackendParams) (*notifymentions.Backend, error) {
	delivery, err := createDelivery(params.client, params.serverRoot)
	if err != nil {
		return nil, err
	}

	backend := notifymentions.New(delivery, params.logger)

	return backend, nil
}

func createSubscriptionsNotifyBackend(params notifyBackendParams, store store.Store) (*notifysubscriptions.Backend, error) {
	delivery, err := createDelivery(params.client, params.serverRoot)
	if err != nil {
		return nil, err
	}

	backend := notifysubscriptions.New(params.serverRoot, store, delivery, params.logger)

	return backend, nil
}

func createDelivery(client *pluginapi.Client, serverRoot string) (*plugindelivery.PluginDelivery, error) {
	bot := &model.Bot{
		Username:    botUsername,
		DisplayName: botDisplayname,
		Description: botDescription,
	}
	botID, err := client.Bot.EnsureBot(bot)
	if err != nil {
		return nil, fmt.Errorf("failed to ensure %s bot: %w", botDisplayname, err)
	}

	pluginAPI := &pluginAPIAdapter{client: client}

	return plugindelivery.New(botID, serverRoot, pluginAPI), nil
}

type pluginAPIAdapter struct {
	client *pluginapi.Client
}

func (da *pluginAPIAdapter) GetDirectChannel(userID1, userID2 string) (*model.Channel, error) {
	return da.client.Channel.GetDirect(userID1, userID2)
}

func (da *pluginAPIAdapter) CreatePost(post *model.Post) error {
	return da.client.Post.CreatePost(post)
}

func (da *pluginAPIAdapter) GetUserByID(userID string) (*model.User, error) {
	return da.client.User.Get(userID)
}

func (da *pluginAPIAdapter) GetUserByUsername(name string) (*model.User, error) {
	return da.client.User.GetByUsername(name)
}

func (da *pluginAPIAdapter) GetTeamMember(teamID string, userID string) (*model.TeamMember, error) {
	return da.client.Team.GetMember(teamID, userID)
}

func (da *pluginAPIAdapter) GetChannelByID(channelID string) (*model.Channel, error) {
	return da.client.Channel.Get(channelID)
}