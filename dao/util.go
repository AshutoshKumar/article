package dao

//Insert in Mongodb
func (a *App) Insert(article Articles) error {
	err := a.MongodbConnection.Insert(&article)

	if err != nil {
		return err
	}
	return nil
}

func (a *App) Delete(movie Movie) error {
	err := db.C(COLLECTION).Remove(&movie)
	return err
}

func (a *App) Update(movie Movie) error {
	err := db.C(COLLECTION).UpdateId(movie.ID, &movie)
	return err
}
