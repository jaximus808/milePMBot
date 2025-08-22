package util

import (
	"encoding/json"

	"github.com/jaximus808/milePMBot/internal/supabaseutil"
)

func (s SupaDB) DBGetUserProfilesExists(discordIds []string) (*[]UserProfile, error) {
	var userProfiles []UserProfile
	res, _, err := supabaseutil.Client.From("UserProfile").Select("*", "", false).In("discord_id", discordIds).Execute()
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(res, &userProfiles)
	if err != nil {
		return nil, err
	}

	return &userProfiles, nil
}
