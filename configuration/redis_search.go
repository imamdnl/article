package configuration

import (
	"github.com/RediSearch/redisearch-go/redisearch"
	"log"
)

func ConfigSearch(addr string) *redisearch.Client {
	// Create a client. By default a client is schemaless
	// unless a schema is provided when creating the index
	c := redisearch.NewClient(addr, "myIndex")

	// Create a schema
	sc := redisearch.NewSchema(redisearch.DefaultOptions).
		AddField(redisearch.NewTextField("Id")).
		AddField(redisearch.NewTextField("Author")).
		AddField(redisearch.NewTextField("Title")).
		AddField(redisearch.NewTextFieldOptions("Body", redisearch.TextFieldOptions{Weight: 5.0})).
		AddField(redisearch.NewTextFieldOptions("Created", redisearch.TextFieldOptions{Sortable: true}))

	// Drop an existing index. If the index does not exist an error is returned
	c.Drop()

	// Create the index with the given schema
	if err := c.CreateIndex(sc); err != nil {
		log.Fatal(err)
	}
	return c
}
