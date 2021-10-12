# Link Shortener
Link Shortener in Go, using Redis as a database

## Available routes

	GET /{id} -> redirect to the original url

	POST /create -> create a link
        {
			"ID": "this is my link id",
			"Location": "https://www.youtube.com/watch?v=dQw4w9WgXcQ",
			"Password": "this is used to delete your link"
		}

	DELETE /delete -> delete a link

        {
			"ID": "id of the link to be deleted",
			"Password": "so only you can delete it!"
		}
	GET /clicks/{id} -> number of times a link was clicked