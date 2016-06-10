// khan
// https://github.com/topfreegames/khan
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2016 Top Free Games <backend@tfgco.com>

package models

import (
	"fmt"
	"testing"

	"github.com/Pallinder/go-randomdata"
	. "github.com/franela/goblin"
)

func TestClanModel(t *testing.T) {
	g := Goblin(t)
	testDb, err := GetTestDB()
	g.Assert(err == nil).IsTrue()

	g.Describe("Clan Model", func() {
		g.It("Should create a new Clan", func() {
			player := PlayerFactory.MustCreate().(*Player)
			err := testDb.Insert(player)
			g.Assert(err == nil).IsTrue()

			clan := &Clan{
				GameID:   "test",
				PublicID: "test-clan-2",
				Name:     "clan-name",
				Metadata: "{}",
				OwnerID:  player.ID,
			}
			err = testDb.Insert(clan)
			g.Assert(err == nil).IsTrue()
			g.Assert(clan.ID != 0).IsTrue()

			dbClan, err := GetClanByID(clan.ID)
			g.Assert(err == nil).IsTrue()

			g.Assert(dbClan.GameID).Equal(clan.GameID)
			g.Assert(dbClan.PublicID).Equal(clan.PublicID)
		})

		g.It("Should update a Clan", func() {
			player := PlayerFactory.MustCreate().(*Player)
			err := testDb.Insert(player)
			g.Assert(err == nil).IsTrue()

			clan := ClanFactory.MustCreateWithOption(map[string]interface{}{
				"OwnerID": player.ID,
			}).(*Clan)
			err = testDb.Insert(clan)
			g.Assert(err == nil).IsTrue()
			dt := clan.UpdatedAt

			clan.Metadata = "{ \"x\": 1 }"
			count, err := testDb.Update(clan)
			g.Assert(err == nil).IsTrue()
			g.Assert(int(count)).Equal(1)
			g.Assert(clan.UpdatedAt > dt).IsTrue()
		})

		g.It("Should get existing Clan", func() {
			player := PlayerFactory.MustCreate().(*Player)
			err := testDb.Insert(player)
			g.Assert(err == nil).IsTrue()

			clan := ClanFactory.MustCreateWithOption(map[string]interface{}{
				"OwnerID": player.ID,
			}).(*Clan)
			err = testDb.Insert(clan)
			g.Assert(err == nil).IsTrue()

			dbClan, err := GetClanByID(clan.ID)
			g.Assert(err == nil).IsTrue()
			g.Assert(dbClan.ID).Equal(clan.ID)
		})

		g.It("Should not get non-existing Clan", func() {
			_, err := GetClanByID(-1)
			g.Assert(err != nil).IsTrue()
			g.Assert(err.Error()).Equal("Clan was not found with id: -1")
		})

		g.It("Should get an existing Clan by Game and PublicID", func() {
			player := PlayerFactory.MustCreate().(*Player)
			err := db.Insert(player)
			g.Assert(err == nil).IsTrue()

			clan := ClanFactory.MustCreateWithOption(map[string]interface{}{
				"OwnerID": player.ID,
			}).(*Clan)
			err = testDb.Insert(clan)
			g.Assert(err == nil).IsTrue()

			dbClan, err := GetClanByPublicID(clan.GameID, clan.PublicID)
			g.Assert(err == nil).IsTrue()
			g.Assert(dbClan.ID).Equal(clan.ID)
		})

		g.It("Should not get a non-existing Clan by Game and PublicID", func() {
			_, err := GetClanByPublicID("invalid-game", "invalid-clan")
			g.Assert(err != nil).IsTrue()
			g.Assert(err.Error()).Equal("Clan was not found with id: invalid-clan")
		})

		g.It("Should get an existing Clan by Game, PublicID and OwnerPublicID", func() {
			player := PlayerFactory.MustCreate().(*Player)
			err := db.Insert(player)
			g.Assert(err == nil).IsTrue()

			clan := ClanFactory.MustCreateWithOption(map[string]interface{}{
				"OwnerID": player.ID,
			}).(*Clan)
			err = testDb.Insert(clan)
			g.Assert(err == nil).IsTrue()

			dbClan, err := GetClanByPublicIDAndOwnerPublicID(clan.GameID, clan.PublicID, player.PublicID)
			g.Assert(err == nil).IsTrue()
			g.Assert(dbClan.ID).Equal(clan.ID)
			g.Assert(dbClan.GameID).Equal(clan.GameID)
			g.Assert(dbClan.PublicID).Equal(clan.PublicID)
			g.Assert(dbClan.Name).Equal(clan.Name)
			g.Assert(dbClan.OwnerID).Equal(clan.OwnerID)
		})

		g.It("Should not get a non-existing Clan by Game, PublicID and OwnerPublicID", func() {
			_, err := GetClanByPublicIDAndOwnerPublicID("invalid-game", "invalid-clan", "invalid-owner-public-id")
			g.Assert(err != nil).IsTrue()
			g.Assert(err.Error()).Equal("Clan was not found with id: invalid-clan")
		})

		g.It("Should not get a existing Clan by Game, PublicID and OwnerPublicID if not Clan owner", func() {
			player := PlayerFactory.MustCreate().(*Player)
			err := db.Insert(player)
			g.Assert(err == nil).IsTrue()

			clan := ClanFactory.MustCreateWithOption(map[string]interface{}{
				"OwnerID": player.ID,
			}).(*Clan)
			err = testDb.Insert(clan)
			g.Assert(err == nil).IsTrue()

			_, err = GetClanByPublicIDAndOwnerPublicID(clan.GameID, clan.PublicID, "invalid-owner-public-id")
			g.Assert(err != nil).IsTrue()
			g.Assert(err.Error()).Equal(fmt.Sprintf("Clan was not found with id: %s", clan.PublicID))
		})

		g.It("Should create a new Clan with CreateClan", func() {
			player := PlayerFactory.MustCreate().(*Player)
			err := db.Insert(player)
			g.Assert(err == nil).IsTrue()

			clan, err := CreateClan(
				player.GameID,
				randomdata.FullName(randomdata.RandomGender),
				"clan-name",
				player.PublicID,
				"{}",
			)

			g.Assert(err == nil).IsTrue()
			g.Assert(clan.ID != 0).IsTrue()

			dbClan, err := GetClanByID(clan.ID)
			g.Assert(err == nil).IsTrue()

			g.Assert(dbClan.GameID).Equal(clan.GameID)
			g.Assert(dbClan.PublicID).Equal(clan.PublicID)
		})

		g.It("Should not create a new Clan with CreateClan if invalid data", func() {
			player := PlayerFactory.MustCreate().(*Player)
			err := db.Insert(player)
			g.Assert(err == nil).IsTrue()

			_, err = CreateClan(
				player.GameID,
				randomdata.FullName(randomdata.RandomGender),
				"clan-name",
				player.PublicID,
				"it-will-fail-because-metadata-is-not-a-json",
			)

			g.Assert(err != nil).IsTrue()
			g.Assert(err.Error()).Equal("pq: invalid input syntax for type json")
		})

		g.It("Should not create a new Clan with CreateClan if unexistent player", func() {
			playerPublicID := randomdata.FullName(randomdata.RandomGender)
			_, err := CreateClan(
				"create-1",
				randomdata.FullName(randomdata.RandomGender),
				"clan-name",
				playerPublicID,
				"{}",
			)

			g.Assert(err != nil).IsTrue()
			g.Assert(err.Error()).Equal(fmt.Sprintf("Player was not found with id: %s", playerPublicID))
		})

		g.It("Should update a Clan with UpdateClan", func() {
			player := PlayerFactory.MustCreate().(*Player)
			err := testDb.Insert(player)
			g.Assert(err == nil).IsTrue()

			clan := ClanFactory.MustCreateWithOption(map[string]interface{}{
				"OwnerID": player.ID,
			}).(*Clan)
			err = testDb.Insert(clan)
			g.Assert(err == nil).IsTrue()

			metadata := "{\"x\": 1}"
			updClan, err := UpdateClan(
				clan.GameID,
				clan.PublicID,
				clan.Name,
				player.PublicID,
				metadata,
			)

			g.Assert(err == nil).IsTrue()
			g.Assert(updClan.ID).Equal(clan.ID)

			dbClan, err := GetClanByPublicID(clan.GameID, clan.PublicID)
			g.Assert(err == nil).IsTrue()

			g.Assert(dbClan.Metadata).Equal(metadata)
		})

		g.It("Should not update a Clan if player is not the clan owner with UpdateClan", func() {
			player := PlayerFactory.MustCreate().(*Player)
			err := testDb.Insert(player)
			g.Assert(err == nil).IsTrue()

			owner := PlayerFactory.MustCreate().(*Player)
			err = testDb.Insert(owner)
			g.Assert(err == nil).IsTrue()

			clan := ClanFactory.MustCreateWithOption(map[string]interface{}{
				"OwnerID": owner.ID,
			}).(*Clan)
			err = testDb.Insert(clan)
			g.Assert(err == nil).IsTrue()

			metadata := "{\"x\": 1}"
			_, err = UpdateClan(
				clan.GameID,
				clan.PublicID,
				clan.Name,
				player.PublicID,
				metadata,
			)

			g.Assert(err == nil).IsFalse()
			g.Assert(err.Error()).Equal(fmt.Sprintf("Clan was not found with id: %s", clan.PublicID))
		})

		g.It("Should not update a Clan with Invalid Data with UpdateClan", func() {
			player := PlayerFactory.MustCreate().(*Player)
			err := testDb.Insert(player)
			g.Assert(err == nil).IsTrue()

			clan := ClanFactory.MustCreateWithOption(map[string]interface{}{
				"OwnerID": player.ID,
			}).(*Clan)
			err = testDb.Insert(clan)
			g.Assert(err == nil).IsTrue()

			metadata := "it will not work because i am not a json"
			_, err = UpdateClan(
				clan.GameID,
				clan.PublicID,
				clan.Name,
				player.PublicID,
				metadata,
			)

			g.Assert(err == nil).IsFalse()
			g.Assert(err.Error()).Equal("pq: invalid input syntax for type json")
		})
	})
}
