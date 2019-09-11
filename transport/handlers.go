package transport

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"savingDeposits"
)

func getUsersHandler(service savingDeposits.UserService) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		var input savingDeposits.UserReadInput

		input.Id = vars["id"]
		result, err := service.Read(input)
		if err != nil {
			badRequestError(err, w)
			return
		}

		respond(w, http.StatusOK, result)
	}
}

func getAllUsersHandler(service savingDeposits.UserService) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var input savingDeposits.UserAllInput

		result, err := service.All(input)
		if err != nil {
			respond(w, http.StatusBadRequest, err.Error())
			return
		}

		respond(w, http.StatusOK, result)
	}
}

func postUsersHandler(service savingDeposits.UserService) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		var newUser savingDeposits.UserCreateInput
		if err := json.NewDecoder(r.Body).Decode(&newUser); err != nil {
			respond(w, http.StatusBadRequest, err.Error())
			return
		}

		result, err := service.Create(newUser)
		if err != nil {
			badRequestError(err, w)
			return
		}

		respond(w, http.StatusCreated, result)
	}
}

func patchUsersHandler(service savingDeposits.UserService) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		defer r.Body.Close()

		var input savingDeposits.UserUpdateInput

		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			respond(w, http.StatusBadRequest, err.Error())
			log.Println(err.Error())
			return
		}

		input.Id = vars["id"]
		result, err := service.Update(input)
		if err != nil {
			badRequestError(err, w)
			return
		}

		respond(w, http.StatusOK, result)
	}
}

func deleteUsersHandler(service savingDeposits.UserService) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		var deleteIn savingDeposits.UserDeleteInput
		deleteIn.Id = vars["id"]

		_, err := service.Delete(deleteIn)
		if err != nil {
			badRequestError(err, w)
			return
		}

		respond(w, http.StatusNoContent, nil)
	}
}

func getDepositsHandler(srv savingDeposits.DepositsService) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var input savingDeposits.DepositReadInput
		vars := mux.Vars(r)
		input.Id = vars["id"]
		user, err := tryGetUserFromContext(r.Context())
		if err != nil {
			respond(w, http.StatusUnauthorized, err.Error())
			return
		}
		input.User = *user

		result, err := srv.Read(input)
		if err != nil {
			badRequestError(err, w)
			return
		}

		respond(w, http.StatusOK, result)
	}
}

func getAllDepositsHandler(srv savingDeposits.DepositsService) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var input savingDeposits.DepositFindInput
		input.Query = r.URL.RawQuery
		user, err := tryGetUserFromContext(r.Context())
		if err != nil {
			respond(w, http.StatusUnauthorized, err.Error())
			return
		}

		input.User = *user
		result, err := srv.Find(input)
		if err != nil {
			respond(w, http.StatusBadRequest, err.Error())
			return
		}

		respond(w, http.StatusOK, result)
	}
}

func postDepositsHandler(srv savingDeposits.DepositsService) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		var input savingDeposits.DepositCreateInput

		if err := json.NewDecoder(r.Body).Decode(&input.SavingDeposit); err != nil {
			respond(w, http.StatusBadRequest, err.Error())
			return
		}

		user, err := tryGetUserFromContext(r.Context())
		if err != nil {
			respond(w, http.StatusUnauthorized, err.Error())
			return
		}
		input.User = *user

		result, err := srv.Create(input)
		if err != nil {
			badRequestError(err, w)
			return
		}

		respond(w, http.StatusCreated, result)
	}
}

func patchDepositsHandler(srv savingDeposits.DepositsService) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		var input savingDeposits.DepositUpdateInput
		vars := mux.Vars(r)
		input.Id = vars["id"]

		if err := json.NewDecoder(r.Body).Decode(&input.Data); err != nil {
			respond(w, http.StatusBadRequest, err.Error())
			return
		}

		user, err := tryGetUserFromContext(r.Context())
		if err != nil {
			respond(w, http.StatusForbidden, err.Error())
			return
		}
		input.User = *user

		result, err := srv.Update(input)
		if err != nil {
			badRequestError(err, w)
			return
		}

		respond(w, http.StatusOK, result)
	}
}

func deleteDepositsHandler(srv savingDeposits.DepositsService) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var input savingDeposits.DepositDeleteInput

		user, err := tryGetUserFromContext(r.Context())
		if err != nil {
			respond(w, http.StatusUnauthorized, err.Error())
			return
		}

		vars := mux.Vars(r)
		input.Id = vars["id"]
		input.User = *user

		_, err = srv.Delete(input)
		if err != nil {
			if err == savingDeposits.NotAuthorizedError {
				respond(w, http.StatusForbidden, err.Error())
				return
			}

			badRequestError(err, w)
			return
		}

		respond(w, http.StatusNoContent, nil)
	}
}

func (s *Server) loginHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var userData struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}

		defer r.Body.Close()
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&userData)
		if err != nil {
			respond(w, http.StatusInternalServerError, "Internal Server error")
			log.Printf("[ERROR] %v", err)
			return
		}

		token, err := s.authn.Login(userData.Username, userData.Password)
		if err != nil {
			respond(w, http.StatusUnauthorized, "Not allowed")
			return
		}

		var returnToken struct {
			Token string `json:"token"`
		}
		returnToken.Token = token

		respond(w, http.StatusOK, returnToken)
	}
}

func (s *Server) profileHandler() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, err := tryGetUserFromContext(r.Context())
		if err != nil {
			respond(w, http.StatusUnauthorized, err.Error())
			return
		}

		respond(w, http.StatusOK, user)
	})
}

func (s *Server) newClientHandler() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var newClient struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}

		defer r.Body.Close()
		if err := json.NewDecoder(r.Body).Decode(&newClient); err != nil {
			respond(w, http.StatusInternalServerError, "Internal Server error")
			log.Printf("[ERROR] %v", err)
			return
		}

		user, err := s.userService.Create(savingDeposits.UserCreateInput{
			Username: newClient.Username,
			Password: newClient.Password,
			Role:     "regular",
		})

		if err != nil {
			respond(w, http.StatusInternalServerError, err.Error())
			log.Printf("[ERROR] %v", err)
			return
		}

		respond(w, http.StatusCreated, user)
	})
}

func (s *Server) generateReportHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var input savingDeposits.GenerateReportInput
		input.Query = r.URL.RawQuery
		user, err := tryGetUserFromContext(r.Context())
		if err != nil {
			respond(w, http.StatusUnauthorized, err.Error())
			return
		}

		input.User = *user
		result, err := s.depositService.GenerateReport(input)
		if err != nil {
			respond(w, http.StatusBadRequest, err.Error())
			return
		}

		respond(w, http.StatusOK, result)
	}
}

func badRequestError(err error, w http.ResponseWriter) {
	log.Printf("[ERROR] %s", err.Error())
	switch err {
	case savingDeposits.NotFoundError:
		respond(w, http.StatusNotFound, err.Error())
	default:
		respond(w, http.StatusBadRequest, err.Error())
	}
}

func tryGetUserFromContext(ctx context.Context) (*savingDeposits.User, error) {
	userRaw := ctx.Value("authUser")
	if userRaw == nil {
		return nil, errors.New("user must be authenticated")
	}

	user, ok := userRaw.(*savingDeposits.User)
	if !ok {
		return nil, errors.New("error fetching authenticated user")
	}
	return user, nil
}
