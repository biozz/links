package migrations

import (
	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

func init() {
	m.Register(func(app core.App) error {
		collection, err := app.FindCollectionByNameOrId("items")
		if err != nil {
			return err
		}
		collection.Fields.Add(&core.JSONField{Name: "default_args"})
		if err := app.Save(collection); err != nil {
			return err
		}
		_, err = app.DB().NewQuery("UPDATE items SET default_args = '[]'").Execute()
		return err
	}, func(app core.App) error {
		collection, err := app.FindCollectionByNameOrId("items")
		if err != nil {
			return err
		}
		collection.Fields.RemoveByName("default_args")
		return app.Save(collection)
	})
}
