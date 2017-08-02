package phonesearch

type Server *martini.ClassicMartini

func NewServer(session *DatabaseSession) Server {
	// Create the server
	m := Server(martini.Classic())
	m.Use(render.Renderer(render.Options{
		IndentJSON: true,
	}))
	m.Use(session.Database())

	// Define the "GET /signatures" route.
	m.Get("/reverse", func(r render.Render, db *mgo.Database) {
		r.JSON(200, fetchAllSignatures(db))
	})

	return m
}

func fetchPersonByNumber(phoneNumber string) {
	// TODO look up in DB first if not found call scraper

}