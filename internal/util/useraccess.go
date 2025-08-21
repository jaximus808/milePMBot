package util

import (
	"encoding/json"

	"github.com/jaximus808/milePMBot/internal/supabaseutil"
)

func (s SupaDB) DBGetUserAccessExists(discordIds []string) (*[]UserAccess, error) {
	var userAccess []UserAccess
	res, _, err := supabaseutil.Client.From("UserAccess").Select("*", "", false).In("discord_id", discordIds).Execute()
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(res, &userAccess)
	if err != nil {
		return nil, err
	}

	return &userAccess, nil
}
