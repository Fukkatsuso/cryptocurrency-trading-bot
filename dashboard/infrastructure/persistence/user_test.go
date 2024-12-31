package persistence_test

// func TestUser(t *testing.T) {
// 	db := persistence.NewMySQLTransaction(config.DSN())
// 	defer db.Rollback()

// 	userRepository := persistence.NewUserRepository(db)

// 	t.Run("save", func(t *testing.T) {
// 		user := model.NewUser("test", "password", "")
// 		err := userRepository.Save(user)
// 		if err != nil {
// 			t.Fatal(err.Error())
// 		}
// 	})

// 	t.Run("find by user_id", func(t *testing.T) {
// 		user, err := userRepository.FindByID("test")
// 		if err != nil {
// 			t.Fatal(err.Error())
// 		}
// 		if user.Password() != "password" {
// 			t.Fatalf("%s != %s", user.Password(), "password")
// 		}
// 	})

// 	t.Run("update password", func(t *testing.T) {
// 		// "test"ユーザのパスワードを"qwerty"に更新
// 		user := model.NewUser("test", "qwerty", "")
// 		err := userRepository.Save(user)
// 		if err != nil {
// 			t.Fatal(err.Error())
// 		}

// 		user, err = userRepository.FindByID("test")
// 		if err != nil {
// 			t.Fatal(err.Error())
// 		}
// 		if user.Password() != "qwerty" {
// 			t.Fatalf("%s != %s", user.Password(), "qwerty")
// 		}
// 	})
// }
