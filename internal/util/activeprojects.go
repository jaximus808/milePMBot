package util

import (
	"encoding/json"
	"strconv"

	"github.com/jaximus808/milePMBot/internal/supabaseutil"
)

/**
Active Project Methods
*/

func (s SupaDB) DBCreateActiveProject(guildId int, pchannelId int, projectId int) (*ActiveProject, error) {
	var selectedActiveProject ActiveProject
	insertedActiveProject := ActiveProjectInsert{
		GuildID:    &guildId,
		PChannelID: &pchannelId,
		ProjectID:  &projectId,
	}
	res, _, err := supabaseutil.Client.From("ActiveProjects").Insert(insertedActiveProject, false, "", "representation", "").Single().Execute()
	if err != nil {

		return nil, err
	}
	err = json.Unmarshal(res, &selectedActiveProject)
	if err != nil {

		return nil, err
	}
	return &selectedActiveProject, nil
}

func (s SupaDB) DBGetActiveProject(guildId string, pchannelId string) (*ActiveProject, error) {
	var selectedActiveProject ActiveProject

	res, _, err := supabaseutil.Client.From("ActiveProjects").Select("*", "", false).Eq("pChannelId", pchannelId).Eq("guildId", guildId).Single().Execute()
	if err != nil {

		return nil, err
	}
	err = json.Unmarshal(res, &selectedActiveProject)
	if err != nil {

		return nil, err
	}
	return &selectedActiveProject, nil
}

func (s SupaDB) DBUpdateProjectId(projectId int, guildId int, pchannelId int) (*ActiveProject, error) {
	var newActiveProejct ActiveProject
	updatedActiveProject := ActiveProjectUpdate{
		GuildID:    &guildId,
		PChannelID: &pchannelId,
	}

	res, _, err := supabaseutil.Client.From("ActiveProjects").Update(updatedActiveProject, "representation", "").Eq("project_id", strconv.Itoa(projectId)).Single().Execute()
	if err != nil {

		return nil, err
	}
	err = json.Unmarshal(res, &newActiveProejct)
	if err != nil {

		return nil, err
	}
	return &newActiveProejct, nil
}

func (s SupaDB) DBEndActiveProject(projectId int) error {
	_, _, err := supabaseutil.Client.From("ActiveProjects").Delete("*", "").Eq("project_id", strconv.Itoa(projectId)).Execute()
	return err
}
