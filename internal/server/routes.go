package server

func (s *server) SetUpRoutes() {
	s.router.POST("/", s.CreateUser())
	s.router.PATCH("/", s.UpdateUser())
	s.router.DELETE("/:user", s.DeleteUserByID())
	s.router.GET("/:user", s.GetUserByID())

	s.router.POST("/order", s.CreateOrder())
}
