// khan
// https://github.com/jpholanda/khan
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2016 Top Free Games <backend@tfgco.com>

package models

import (
	"database/sql"
	"encoding/json"
)

type clanDetailsDAO struct {
	// Clan general information
	GameID               string
	ClanPublicID         string
	ClanName             string
	ClanMetadata         map[string]interface{}
	ClanAllowApplication bool
	ClanAutoJoin         bool
	ClanMembershipCount  int

	//Membership Information
	MembershipLevel      sql.NullString
	MembershipApproved   sql.NullBool
	MembershipDenied     sql.NullBool
	MembershipBanned     sql.NullBool
	MembershipCreatedAt  sql.NullInt64
	MembershipUpdatedAt  sql.NullInt64
	MembershipApprovedAt sql.NullInt64
	MembershipDeniedAt   sql.NullInt64
	MembershipMessage    sql.NullString

	// Clan Owner Information
	OwnerPublicID string
	OwnerName     string
	OwnerMetadata map[string]interface{}

	// Member Information
	PlayerPublicID   sql.NullString
	PlayerName       sql.NullString
	DBPlayerMetadata sql.NullString
	PlayerMetadata   map[string]interface{}
	MembershipCount  int
	OwnershipCount   int

	// Requestor Information
	RequestorPublicID sql.NullString
	RequestorName     sql.NullString

	// Approver Information
	ApproverPublicID sql.NullString
	ApproverName     sql.NullString

	// Denier Information
	DenierPublicID sql.NullString
	DenierName     sql.NullString
}

func (member *clanDetailsDAO) Serialize(includeMembershipLevel bool) map[string]interface{} {
	result := map[string]interface{}{
		"player": map[string]interface{}{
			"publicID": nullOrString(member.PlayerPublicID),
			"name":     nullOrString(member.PlayerName),
		},
	}
	if member.DBPlayerMetadata.Valid {
		json.Unmarshal([]byte(nullOrString(member.DBPlayerMetadata)), &member.PlayerMetadata)
	} else {
		member.PlayerMetadata = map[string]interface{}{}
	}
	if includeMembershipLevel {
		result["level"] = nullOrString(member.MembershipLevel)
	}
	result["player"].(map[string]interface{})["metadata"] = member.PlayerMetadata

	if member.ApproverName.Valid {
		result["player"].(map[string]interface{})["approver"] = map[string]interface{}{
			"name":     member.ApproverName.String,
			"publicID": member.ApproverPublicID.String,
		}
	} else if member.DenierName.Valid {
		result["player"].(map[string]interface{})["denier"] = map[string]interface{}{
			"name":     member.DenierName.String,
			"publicID": member.DenierPublicID.String,
		}
	}
	return result
}

type playerDetailsDAO struct {
	// Player Details
	PlayerID        int
	PlayerName      string
	PlayerMetadata  map[string]interface{}
	PlayerPublicID  string
	PlayerCreatedAt int64
	PlayerUpdatedAt int64

	// Membership Details
	MembershipLevel      sql.NullString
	MembershipApproved   sql.NullBool
	MembershipDenied     sql.NullBool
	MembershipBanned     sql.NullBool
	MembershipCreatedAt  sql.NullInt64
	MembershipUpdatedAt  sql.NullInt64
	MembershipDeletedAt  sql.NullInt64
	MembershipApprovedAt sql.NullInt64
	MembershipDeniedAt   sql.NullInt64
	MembershipMessage    sql.NullString

	// Clan Details
	ClanPublicID        sql.NullString
	ClanName            sql.NullString
	DBClanMetadata      sql.NullString
	ClanMetadata        map[string]interface{}
	ClanOwnerID         sql.NullInt64
	ClanMembershipCount sql.NullInt64

	// Membership Requestor Details
	RequestorName            sql.NullString
	RequestorPublicID        sql.NullString
	DBRequestorMetadata      sql.NullString
	RequestorMetadata        map[string]interface{}
	RequestorMembershipLevel sql.NullString

	// Membership Approver Details
	ApproverName       sql.NullString
	ApproverPublicID   sql.NullString
	DBApproverMetadata sql.NullString
	ApproverMetadata   map[string]interface{}

	// Membership Denier Details
	DenierName       sql.NullString
	DenierPublicID   sql.NullString
	DBDenierMetadata sql.NullString
	DenierMetadata   map[string]interface{}

	// Deleted by Details
	DeletedByName     sql.NullString
	DeletedByPublicID sql.NullString
}

func (p *playerDetailsDAO) Serialize() map[string]interface{} {
	result := map[string]interface{}{
		"level":      nullOrString(p.MembershipLevel),
		"approved":   nullOrBool(p.MembershipApproved),
		"denied":     nullOrBool(p.MembershipDenied),
		"banned":     nullOrBool(p.MembershipBanned),
		"createdAt":  nullOrInt(p.MembershipCreatedAt),
		"updatedAt":  nullOrInt(p.MembershipUpdatedAt),
		"deletedAt":  nullOrInt(p.MembershipDeletedAt),
		"approvedAt": nullOrInt(p.MembershipApprovedAt),
		"deniedAt":   nullOrInt(p.MembershipDeniedAt),
		"message":    nullOrString(p.MembershipMessage),
		"clan": map[string]interface{}{
			"publicID":        nullOrString(p.ClanPublicID),
			"name":            nullOrString(p.ClanName),
			"membershipCount": nullOrInt(p.ClanMembershipCount),
		},
		"requestor": map[string]interface{}{
			"publicID": nullOrString(p.RequestorPublicID),
			"name":     nullOrString(p.RequestorName),
			"level":    nullOrString(p.RequestorMembershipLevel),
		},
	}

	if p.DBClanMetadata.Valid {
		json.Unmarshal([]byte(nullOrString(p.DBClanMetadata)), &p.ClanMetadata)
	} else {
		p.ClanMetadata = map[string]interface{}{}
	}
	result["clan"].(map[string]interface{})["metadata"] = p.ClanMetadata

	if p.DBRequestorMetadata.Valid {
		json.Unmarshal([]byte(nullOrString(p.DBRequestorMetadata)), &p.RequestorMetadata)
	} else {
		p.RequestorMetadata = map[string]interface{}{}
	}
	result["requestor"].(map[string]interface{})["metadata"] = p.RequestorMetadata

	if p.DeletedByPublicID.Valid {
		result["deletedBy"] = map[string]interface{}{
			"publicID": nullOrString(p.DeletedByPublicID),
			"name":     nullOrString(p.DeletedByName),
		}
	}

	if p.ApproverPublicID.Valid {
		if p.DBApproverMetadata.Valid {
			json.Unmarshal([]byte(nullOrString(p.DBApproverMetadata)), &p.ApproverMetadata)
		} else {
			p.ApproverMetadata = map[string]interface{}{}
		}

		result["approver"] = map[string]interface{}{
			"publicID": nullOrString(p.ApproverPublicID),
			"name":     nullOrString(p.ApproverName),
			"metadata": p.ApproverMetadata,
		}
	}

	if p.DenierPublicID.Valid {
		if p.DBDenierMetadata.Valid {
			json.Unmarshal([]byte(nullOrString(p.DBDenierMetadata)), &p.DenierMetadata)
		} else {
			p.DenierMetadata = map[string]interface{}{}
		}

		result["denier"] = map[string]interface{}{
			"publicID": nullOrString(p.DenierPublicID),
			"name":     nullOrString(p.DenierName),
			"metadata": p.DenierMetadata,
		}
	}
	return result
}
