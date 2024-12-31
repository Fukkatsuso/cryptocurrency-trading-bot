package service_test

// func TestAuth(t *testing.T) {
// 	tx := persistence.NewMySQLTransaction(config.DSN())
// 	defer tx.Rollback()

// 	userRepository := persistence.NewUserRepository(tx)
// 	sessionRepository := persistence.NewSessionRepository(tx)

// 	authService := service.NewAuthService(userRepository, sessionRepository)

// 	// create testUser
// 	passwordHash, err := model.PasswordHash("password")
// 	if err != nil {
// 		t.Fatal(err.Error())
// 	}
// 	testUser := model.NewUser("test", passwordHash, "")
// 	if err := userRepository.Save(testUser); err != nil {
// 		t.Fatal(err.Error())
// 	}

// 	t.Run("succeed in login", func(t *testing.T) {
// 		sessID, err := authService.Login("test", "password")
// 		if err != nil {
// 			t.Fatal(err.Error())
// 		}
// 		if sessID == "" {
// 			t.Fatal("Login() returns empty sessionID")
// 		}

// 		if !authService.LoggedIn("test", sessID) {
// 			t.Fatal("test user is not logged in")
// 		}

// 		err = authService.Logout("test")
// 		if err != nil {
// 			t.Fatal(err.Error())
// 		}
// 	})

// 	t.Run("fail to login by wrong id", func(t *testing.T) {
// 		sessID, err := authService.Login("testtest", "password")
// 		if err == nil {
// 			t.Fatal("Login() must fail")
// 		}

// 		if authService.LoggedIn("test", sessID) {
// 			t.Fatal("LoggedId() must return false")
// 		}
// 	})

// 	t.Run("fail to login by wrong password", func(t *testing.T) {
// 		sessID, err := authService.Login("test", "passwordpassword")
// 		if err == nil {
// 			t.Fatal("Login() must fail")
// 		}

// 		if authService.LoggedIn("test", sessID) {
// 			t.Fatal("LoggedId() must return false")
// 		}
// 	})
// }
