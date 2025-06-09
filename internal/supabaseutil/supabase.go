package supabaseutil

import (
	"fmt"
	"os"

	"github.com/supabase-community/supabase-go"
)

var Client *supabase.Client

func InitializeSupabase() error {

	client, err := supabase.NewClient(os.Getenv("APP_SUPABASE_URL"), os.Getenv("APP_SUPABASE_ANON_KEY"), &supabase.ClientOptions{})
	if err != nil {
		return fmt.Errorf("failed to open supabase client: %w", err)
	}
	Client = client
	return nil

}
