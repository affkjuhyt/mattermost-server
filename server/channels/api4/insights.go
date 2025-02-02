// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package api4

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/mattermost/mattermost-server/server/public/model"
)

func (api *API) InitInsights() {
	// Reactions
	api.BaseRoutes.InsightsForTeam.Handle("/reactions", api.APISessionRequired(getTopReactionsForTeamSince)).Methods("GET")
	api.BaseRoutes.InsightsForUser.Handle("/reactions", api.APISessionRequired(getTopReactionsForUserSince)).Methods("GET")

	// Channels
	api.BaseRoutes.InsightsForTeam.Handle("/channels", api.APISessionRequired(getTopChannelsForTeamSince)).Methods("GET")
	api.BaseRoutes.InsightsForUser.Handle("/channels", api.APISessionRequired(getTopChannelsForUserSince)).Methods("GET")

	// Threads
	api.BaseRoutes.InsightsForTeam.Handle("/threads", api.APISessionRequired(getTopThreadsForTeamSince)).Methods("GET")
	api.BaseRoutes.InsightsForUser.Handle("/threads", api.APISessionRequired(getTopThreadsForUserSince)).Methods("GET")

	// user DMs
	api.BaseRoutes.InsightsForUser.Handle("/dms", api.APISessionRequired(getTopDMsForUserSince)).Methods("GET")

	// Inactive channels
	api.BaseRoutes.InsightsForTeam.Handle("/inactive_channels", api.APISessionRequired(getTopInactiveChannelsForTeamSince)).Methods("GET")
	api.BaseRoutes.InsightsForUser.Handle("/inactive_channels", api.APISessionRequired(getTopInactiveChannelsForUserSince)).Methods("GET")

	// New teammembers
	api.BaseRoutes.InsightsForTeam.Handle("/team_members", api.APISessionRequired(getNewTeamMembersSince)).Methods("GET")
}

// Top Reactions

func getTopReactionsForTeamSince(c *Context, w http.ResponseWriter, r *http.Request) {

	// license and guest user check
	permissionErr := minimumProfessionalLicense(c)
	if permissionErr != nil {
		c.Err = permissionErr
		return
	}
	permissionErr = rejectGuests(c)
	if permissionErr != nil {
		c.Err = permissionErr
		return
	}

	c.RequireTeamId()
	if c.Err != nil {
		return
	}

	team, appErr := c.App.GetTeam(c.Params.TeamId)
	if appErr != nil {
		c.Err = appErr
		return
	}

	if !c.App.SessionHasPermissionToTeam(*c.AppContext.Session(), team.Id, model.PermissionViewTeam) {
		c.SetPermissionError(model.PermissionViewTeam)
		return
	}

	user, appErr := c.App.GetUser(c.AppContext.Session().UserId)
	if appErr != nil {
		c.Err = appErr
		return
	}

	startTime, appErr := model.GetStartOfDayForTimeRange(c.Params.TimeRange, user.GetTimezoneLocation())
	if appErr != nil {
		c.Err = appErr
		return
	}

	topReactionList, appErr := c.App.GetTopReactionsForTeamSince(c.Params.TeamId, c.AppContext.Session().UserId, &model.InsightsOpts{
		StartUnixMilli: startTime.UnixMilli(),
		Page:           c.Params.Page,
		PerPage:        c.Params.PerPage,
	})
	if appErr != nil {
		c.Err = appErr
		return
	}

	if err := json.NewEncoder(w).Encode(topReactionList); err != nil {
		c.Err = model.NewAppError("getTopReactionsForTeamSince", "api.marshal_error", nil, err.Error(), http.StatusInternalServerError)
		return
	}
}

func getTopReactionsForUserSince(c *Context, w http.ResponseWriter, r *http.Request) {
	// guest user check
	permissionErr := rejectGuests(c)
	if permissionErr != nil {
		c.Err = permissionErr
		return
	}

	c.Params.TeamId = r.URL.Query().Get("team_id")

	// TeamId is an optional parameter
	if c.Params.TeamId != "" {
		if !model.IsValidId(c.Params.TeamId) {
			c.SetInvalidURLParam("team_id")
			return
		}

		team, appErr := c.App.GetTeam(c.Params.TeamId)
		if appErr != nil {
			c.Err = appErr
			return
		}

		if !c.App.SessionHasPermissionToTeam(*c.AppContext.Session(), team.Id, model.PermissionViewTeam) {
			c.SetPermissionError(model.PermissionViewTeam)
			return
		}
	}

	user, appErr := c.App.GetUser(c.AppContext.Session().UserId)
	if appErr != nil {
		c.Err = appErr
		return
	}

	startTime, appErr := model.GetStartOfDayForTimeRange(c.Params.TimeRange, user.GetTimezoneLocation())
	if appErr != nil {
		c.Err = appErr
		return
	}

	topReactionList, appErr := c.App.GetTopReactionsForUserSince(c.AppContext.Session().UserId, c.Params.TeamId, &model.InsightsOpts{
		StartUnixMilli: startTime.UnixMilli(),
		Page:           c.Params.Page,
		PerPage:        c.Params.PerPage,
	})
	if appErr != nil {
		c.Err = appErr
		return
	}

	if err := json.NewEncoder(w).Encode(topReactionList); err != nil {
		c.Err = model.NewAppError("getTopReactionsForUserSince", "api.marshal_error", nil, err.Error(), http.StatusInternalServerError)
		return
	}
}

// Top Channels

func getTopChannelsForTeamSince(c *Context, w http.ResponseWriter, r *http.Request) {
	// license and guest user check
	permissionErr := minimumProfessionalLicense(c)
	if permissionErr != nil {
		c.Err = permissionErr
		return
	}
	permissionErr = rejectGuests(c)
	if permissionErr != nil {
		c.Err = permissionErr
		return
	}

	c.RequireTeamId()
	if c.Err != nil {
		return
	}

	team, appErr := c.App.GetTeam(c.Params.TeamId)
	if appErr != nil {
		c.Err = appErr
		return
	}

	if !c.App.SessionHasPermissionToTeam(*c.AppContext.Session(), team.Id, model.PermissionViewTeam) {
		c.SetPermissionError(model.PermissionViewTeam)
		return
	}

	user, appErr := c.App.GetUser(c.AppContext.Session().UserId)
	if appErr != nil {
		c.Err = appErr
		return
	}

	loc := user.GetTimezoneLocation()
	startTime, appErr := model.GetStartOfDayForTimeRange(c.Params.TimeRange, loc)
	if appErr != nil {
		c.Err = appErr
		return
	}

	topChannels, appErr := c.App.GetTopChannelsForTeamSince(c.AppContext, c.Params.TeamId, c.AppContext.Session().UserId, &model.InsightsOpts{
		StartUnixMilli: startTime.UnixMilli(),
		Page:           c.Params.Page,
		PerPage:        c.Params.PerPage,
	})
	if appErr != nil {
		c.Err = appErr
		return
	}

	topChannels.PostCountByDuration, appErr = postCountByDurationViewModel(c, topChannels, startTime, c.Params.TimeRange, nil, loc)
	if appErr != nil {
		c.Err = appErr
		return
	}

	if err := json.NewEncoder(w).Encode(topChannels); err != nil {
		c.Err = model.NewAppError("getTopChannelsForTeamSince", "api.marshal_error", nil, err.Error(), http.StatusInternalServerError)
		return
	}
}

func getTopChannelsForUserSince(c *Context, w http.ResponseWriter, r *http.Request) {
	// guest user check
	permissionErr := rejectGuests(c)
	if permissionErr != nil {
		c.Err = permissionErr
		return
	}

	c.Params.TeamId = r.URL.Query().Get("team_id")

	// TeamId is an optional parameter
	if c.Params.TeamId != "" {
		if !model.IsValidId(c.Params.TeamId) {
			c.SetInvalidURLParam("team_id")
			return
		}

		team, appErr := c.App.GetTeam(c.Params.TeamId)
		if appErr != nil {
			c.Err = appErr
			return
		}

		if !c.App.SessionHasPermissionToTeam(*c.AppContext.Session(), team.Id, model.PermissionViewTeam) {
			c.SetPermissionError(model.PermissionViewTeam)
			return
		}
	}

	user, appErr := c.App.GetUser(c.AppContext.Session().UserId)
	if appErr != nil {
		c.Err = appErr
		return
	}

	loc := user.GetTimezoneLocation()
	startTime, appErr := model.GetStartOfDayForTimeRange(c.Params.TimeRange, loc)
	if appErr != nil {
		c.Err = appErr
		return
	}

	topChannels, appErr := c.App.GetTopChannelsForUserSince(c.AppContext, c.AppContext.Session().UserId, c.Params.TeamId, &model.InsightsOpts{
		StartUnixMilli: startTime.UnixMilli(),
		Page:           c.Params.Page,
		PerPage:        c.Params.PerPage,
	})
	if appErr != nil {
		c.Err = appErr
		return
	}

	topChannels.PostCountByDuration, appErr = postCountByDurationViewModel(c, topChannels, startTime, c.Params.TimeRange, &c.AppContext.Session().UserId, loc)
	if appErr != nil {
		c.Err = appErr
		return
	}

	if err := json.NewEncoder(w).Encode(topChannels); err != nil {
		c.Err = model.NewAppError("getTopChannelsForUserSince", "api.marshal_error", nil, err.Error(), http.StatusInternalServerError)
		return
	}
}

// Top Threads
func getTopThreadsForTeamSince(c *Context, w http.ResponseWriter, r *http.Request) {
	// license and guest user check
	permissionErr := minimumProfessionalLicense(c)
	if permissionErr != nil {
		c.Err = permissionErr
		return
	}
	permissionErr = rejectGuests(c)
	if permissionErr != nil {
		c.Err = permissionErr
		return
	}

	c.RequireTeamId()
	if c.Err != nil {
		return
	}

	team, appErr := c.App.GetTeam(c.Params.TeamId)
	if appErr != nil {
		c.Err = appErr
		return
	}

	// restrict users with no access to team
	user, err := c.App.GetUser(c.AppContext.Session().UserId)
	if err != nil {
		c.Err = err
		return
	}

	if !c.App.SessionHasPermissionToTeam(*c.AppContext.Session(), team.Id, model.PermissionViewTeam) {
		c.SetPermissionError(model.PermissionViewTeam)
		return
	}

	startTime, appErr := model.GetStartOfDayForTimeRange(c.Params.TimeRange, user.GetTimezoneLocation())
	if appErr != nil {
		c.Err = appErr
		return
	}

	topThreads, appErr := c.App.GetTopThreadsForTeamSince(c.AppContext, c.Params.TeamId, c.AppContext.Session().UserId, &model.InsightsOpts{
		StartUnixMilli: startTime.UnixMilli(),
		Page:           c.Params.Page,
		PerPage:        c.Params.PerPage,
	})
	if appErr != nil {
		c.Err = appErr
		return
	}

	if err := json.NewEncoder(w).Encode(topThreads); err != nil {
		c.Err = model.NewAppError("getTopThreadsForTeamSince", "api.marshal_error", nil, err.Error(), http.StatusInternalServerError)
		return
	}
}

func getTopThreadsForUserSince(c *Context, w http.ResponseWriter, r *http.Request) {
	// guest user check
	permissionErr := rejectGuests(c)
	if permissionErr != nil {
		c.Err = permissionErr
		return
	}

	c.Params.TeamId = r.URL.Query().Get("team_id")

	// restrict users with no access to team
	user, err := c.App.GetUser(c.AppContext.Session().UserId)
	if err != nil {
		c.Err = err
		return
	}
	// TeamId is an optional parameter
	if c.Params.TeamId != "" {
		if !model.IsValidId(c.Params.TeamId) {
			c.SetInvalidURLParam("team_id")
			return
		}

		team, teamErr := c.App.GetTeam(c.Params.TeamId)
		if teamErr != nil {
			c.Err = teamErr
			return
		}

		if !c.App.SessionHasPermissionToTeam(*c.AppContext.Session(), team.Id, model.PermissionViewTeam) {
			c.SetPermissionError(model.PermissionViewTeam)
			return
		}
	}

	startTime, appErr := model.GetStartOfDayForTimeRange(c.Params.TimeRange, user.GetTimezoneLocation())
	if appErr != nil {
		c.Err = appErr
		return
	}

	topThreads, appErr := c.App.GetTopThreadsForUserSince(c.AppContext, c.Params.TeamId, c.AppContext.Session().UserId, &model.InsightsOpts{
		StartUnixMilli: startTime.UnixMilli(),
		Page:           c.Params.Page,
		PerPage:        c.Params.PerPage,
	})
	if appErr != nil {
		c.Err = appErr
		return
	}

	if err := json.NewEncoder(w).Encode(topThreads); err != nil {
		c.Err = model.NewAppError("getTopThreadsForUserSince", "api.marshal_error", nil, err.Error(), http.StatusInternalServerError)
		return
	}
}

// Top DMs
func getTopDMsForUserSince(c *Context, w http.ResponseWriter, r *http.Request) {
	// guest user check
	permissionErr := rejectGuests(c)
	if permissionErr != nil {
		c.Err = permissionErr
		return
	}

	user, err := c.App.GetUser(c.AppContext.Session().UserId)
	if err != nil {
		c.Err = err
		return
	}

	startTime, appErr := model.GetStartOfDayForTimeRange(c.Params.TimeRange, user.GetTimezoneLocation())
	if appErr != nil {
		c.Err = appErr
		return
	}

	topDMs, err := c.App.GetTopDMsForUserSince(user.Id, &model.InsightsOpts{
		StartUnixMilli: startTime.UnixMilli(),
		Page:           c.Params.Page,
		PerPage:        c.Params.PerPage,
	})

	if err != nil {
		c.Err = err
		return
	}

	if err := json.NewEncoder(w).Encode(topDMs); err != nil {
		c.Err = model.NewAppError("getTopDMsForUserSince", "api.marshal_error", nil, err.Error(), http.StatusInternalServerError)
		return
	}
}

// Top Channels

func getTopInactiveChannelsForTeamSince(c *Context, w http.ResponseWriter, r *http.Request) {
	// license and guest user check
	permissionErr := minimumProfessionalLicense(c)
	if permissionErr != nil {
		c.Err = permissionErr
		return
	}
	permissionErr = rejectGuests(c)
	if permissionErr != nil {
		c.Err = permissionErr
		return
	}

	c.RequireTeamId()
	if c.Err != nil {
		return
	}

	team, err := c.App.GetTeam(c.Params.TeamId)
	if err != nil {
		c.Err = err
		return
	}

	if !c.App.SessionHasPermissionToTeam(*c.AppContext.Session(), team.Id, model.PermissionViewTeam) {
		c.SetPermissionError(model.PermissionViewTeam)
		return
	}

	user, err := c.App.GetUser(c.AppContext.Session().UserId)
	if err != nil {
		c.Err = err
		return
	}

	loc := user.GetTimezoneLocation()
	startTime, appErr := model.GetStartOfDayForTimeRange(c.Params.TimeRange, loc)
	if appErr != nil {
		c.Err = appErr
		return
	}

	topChannels, err := c.App.GetTopInactiveChannelsForTeamSince(c.AppContext, c.Params.TeamId, c.AppContext.Session().UserId, &model.InsightsOpts{
		StartUnixMilli: startTime.UnixMilli(),
		Page:           c.Params.Page,
		PerPage:        c.Params.PerPage,
	})
	if err != nil {
		c.Err = err
		return
	}

	if err := json.NewEncoder(w).Encode(topChannels); err != nil {
		c.Err = model.NewAppError("getTopInactiveChannelsForTeamSince", "api.marshal_error", nil, err.Error(), http.StatusInternalServerError)
		return
	}
}

// top inactive channels

func getTopInactiveChannelsForUserSince(c *Context, w http.ResponseWriter, r *http.Request) {
	// guest user check
	permissionErr := rejectGuests(c)
	if permissionErr != nil {
		c.Err = permissionErr
		return
	}

	c.Params.TeamId = r.URL.Query().Get("team_id")

	// TeamId is an optional parameter
	if c.Params.TeamId != "" {
		if !model.IsValidId(c.Params.TeamId) {
			c.SetInvalidURLParam("team_id")
			return
		}

		team, teamErr := c.App.GetTeam(c.Params.TeamId)
		if teamErr != nil {
			c.Err = teamErr
			return
		}

		if !c.App.SessionHasPermissionToTeam(*c.AppContext.Session(), team.Id, model.PermissionViewTeam) {
			c.SetPermissionError(model.PermissionViewTeam)
			return
		}
	}

	user, err := c.App.GetUser(c.AppContext.Session().UserId)
	if err != nil {
		c.Err = err
		return
	}

	loc := user.GetTimezoneLocation()
	startTime, appErr := model.GetStartOfDayForTimeRange(c.Params.TimeRange, loc)
	if appErr != nil {
		c.Err = appErr
		return
	}

	topChannels, err := c.App.GetTopInactiveChannelsForUserSince(c.AppContext, c.Params.TeamId, c.AppContext.Session().UserId, &model.InsightsOpts{
		StartUnixMilli: startTime.UnixMilli(),
		Page:           c.Params.Page,
		PerPage:        c.Params.PerPage,
	})

	if err != nil {
		c.Err = err
		return
	}

	if err := json.NewEncoder(w).Encode(topChannels); err != nil {
		c.Err = model.NewAppError("getTopInactiveChannelsForUserSince", "api.marshal_error", nil, err.Error(), http.StatusInternalServerError)
		return
	}
}

// postCountByDurationViewModel expects a list of channels that are pre-authorized for the given user to view.
func postCountByDurationViewModel(c *Context, topChannelList *model.TopChannelList, startTime *time.Time, timeRange string, userID *string, location *time.Location) (model.ChannelPostCountByDuration, *model.AppError) {
	if len(topChannelList.Items) == 0 {
		return nil, nil
	}
	var postCountsByDay []*model.DurationPostCount
	channelIDs := topChannelList.ChannelIDs()
	var grouping model.PostCountGrouping
	if timeRange == model.TimeRangeToday {
		grouping = model.PostsByHour
	} else {
		grouping = model.PostsByDay
	}
	postCountsByDay, err := c.App.PostCountsByDuration(c.AppContext, channelIDs, startTime.UnixMilli(), userID, grouping, location)
	if err != nil {
		return nil, err
	}
	return model.ToDailyPostCountViewModel(postCountsByDay, startTime, model.TimeRangeToNumberDays(timeRange), channelIDs), nil
}

func getNewTeamMembersSince(c *Context, w http.ResponseWriter, r *http.Request) {
	// license and guest user check
	permissionErr := minimumProfessionalLicense(c)
	if permissionErr != nil {
		c.Err = permissionErr
		return
	}
	permissionErr = rejectGuests(c)
	if permissionErr != nil {
		c.Err = permissionErr
		return
	}

	c.RequireTeamId()
	if c.Err != nil {
		return
	}

	team, err := c.App.GetTeam(c.Params.TeamId)
	if err != nil {
		c.Err = err
		return
	}

	if !c.App.SessionHasPermissionToTeam(*c.AppContext.Session(), team.Id, model.PermissionViewTeam) {
		c.SetPermissionError(model.PermissionViewTeam)
		return
	}

	user, err := c.App.GetUser(c.AppContext.Session().UserId)
	if err != nil {
		c.Err = err
		return
	}
	loc := user.GetTimezoneLocation()
	startTime, appErr := model.GetStartOfDayForTimeRange(c.Params.TimeRange, loc)
	if appErr != nil {
		c.Err = appErr
		return
	}

	ntms, count, err := c.App.GetNewTeamMembersSince(c.AppContext, c.Params.TeamId, &model.InsightsOpts{
		StartUnixMilli: startTime.UnixMilli(),
		Page:           c.Params.Page,
		PerPage:        c.Params.PerPage,
	})
	if err != nil {
		c.Err = err
		return
	}

	ntms.TotalCount = count

	if err := json.NewEncoder(w).Encode(ntms); err != nil {
		c.Err = model.NewAppError("getNewTeamembersForTeamSince", "api.marshal_error", nil, err.Error(), http.StatusInternalServerError)
		return
	}
}
