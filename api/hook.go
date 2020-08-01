// khan
// https://github.com/jpholanda/khan
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2016 Top Free Games <backend@tfgco.com>

package api

import (
	"net/http"
	"time"

	"github.com/labstack/echo"
	"github.com/jpholanda/khan/log"
	"github.com/jpholanda/khan/models"
	"github.com/uber-go/zap"
)

//CreateHookHandler is the handler responsible for creating new hooks
func CreateHookHandler(app *App) func(c echo.Context) error {
	return func(c echo.Context) error {
		c.Set("route", "CreateHook")
		start := time.Now()
		gameID := c.Param("gameID")

		db := app.Db(c.StdContext())

		l := app.Logger.With(
			zap.String("source", "CreateHookHandler"),
			zap.String("operation", "createHook"),
			zap.String("gameID", gameID),
		)

		var payload HookPayload

		err := WithSegment("payload", c, func() error {
			if err := LoadJSONPayload(&payload, c, l); err != nil {
				log.E(l, "Failed to parse json payload.", func(cm log.CM) {
					cm.Write(zap.Error(err))
				})
				return err
			}

			return nil
		})
		if err != nil {
			return FailWith(http.StatusBadRequest, err.Error(), c)
		}

		var hook *models.Hook
		err = WithSegment("hook-create", c, func() error {
			log.D(l, "Creating hook...")
			hook, err = models.CreateHook(
				db,
				gameID,
				payload.Type,
				payload.HookURL,
			)

			if err != nil {
				log.E(l, "Failed to create the hook.", func(cm log.CM) {
					cm.Write(zap.Error(err))
				})
				return err
			}

			return nil
		})
		if err != nil {
			return FailWith(http.StatusInternalServerError, err.Error(), c)
		}

		log.I(l, "Created hook successfully.", func(cm log.CM) {
			cm.Write(
				zap.String("hookPublicID", hook.PublicID),
				zap.Duration("duration", time.Now().Sub(start)),
			)
		})
		return SucceedWith(map[string]interface{}{
			"publicID": hook.PublicID,
		}, c)
	}
}

// RemoveHookHandler is the handler responsible for removing existing hooks
func RemoveHookHandler(app *App) func(c echo.Context) error {
	return func(c echo.Context) error {
		c.Set("route", "RemoveHook")
		start := time.Now()
		gameID := c.Param("gameID")
		publicID := c.Param("publicID")

		db := app.Db(c.StdContext())

		l := app.Logger.With(
			zap.String("source", "RemoveHookHandler"),
			zap.String("operation", "removeHook"),
			zap.String("gameID", gameID),
			zap.String("hookPublicID", publicID),
		)

		var err error
		err = WithSegment("hook-remove", c, func() error {
			log.D(l, "Removing hook...")
			err = models.RemoveHook(
				db,
				gameID,
				publicID,
			)

			if err != nil {
				log.E(l, "Failed to remove hook.", func(cm log.CM) {
					cm.Write(zap.Error(err))
				})
				return err
			}
			return nil
		})
		if err != nil {
			return FailWith(http.StatusInternalServerError, err.Error(), c)
		}

		log.I(l, "Hook removed successfully.", func(cm log.CM) {
			cm.Write(zap.Duration("duration", time.Now().Sub(start)))
		})
		return SucceedWith(map[string]interface{}{}, c)
	}
}
