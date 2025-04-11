package api

func (s *APIServer) Init(userId int) error {
	// roleStore := roles.NewStore()
	// should run after creating initial user/admin whatever
	// atm, just add roles of 'admin' and 'owner' to the roles table and assign owner to the first user
	// _, err := roleStore.Create(types.Role{Name: "123"})
	// if err != nil {
	// 	// TODO is panic appropriate?
	// 	log.Panic(err)
	// 	return err
	// }

	return nil
}
