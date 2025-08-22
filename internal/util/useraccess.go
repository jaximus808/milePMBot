package util

import "github.com/jaximus808/milePMBot/internal/supabaseutil"

func (s SupaDB) DBInsertUserAccess(userAccessRows []*UserAccessInsert) error {
	_, _, err := supabaseutil.Client.From("UserAccess").Upsert(userAccessRows, "user_id,project_id", "", "").Execute()

	return err
}
