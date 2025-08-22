package util

import "github.com/jaximus808/milePMBot/internal/supabaseutil"

// going ot need to add a command to do the merge
func (s SupaDB) DBInsertPendingAccess(pendingAccessRows []*PendingAccessInsert) error {
	_, _, err := supabaseutil.Client.From("PendingAcces").Upsert(pendingAccessRows, "user_id,project_id", "", "").Execute()

	return err
}
