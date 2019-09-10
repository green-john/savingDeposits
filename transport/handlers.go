package transport

import (
	"encoding/json"
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

		var updateInput savingDeposits.UserUpdateInput

		if err := json.NewDecoder(r.Body).Decode(&updateInput); err != nil {
			respond(w, http.StatusBadRequest, err.Error())
			log.Println(err.Error())
			return
		}

		updateInput.Id = vars["id"]
		result, err := service.Update(updateInput)
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
		vars := mux.Vars(r)
		var input savingDeposits.DepositReadInput
		input.Id = vars["id"]

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
		//input.Query = r.URL.RawQuery

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
		var newDeposits savingDeposits.SavingDeposit
		if err := json.NewDecoder(r.Body).Decode(&newDeposits); err != nil {
			respond(w, http.StatusBadRequest, err.Error())
			return
		}

		result, err := srv.Create(savingDeposits.DespositCreateInput{SavingDeposit: newDeposits})
		if err != nil {
			badRequestError(err, w)
			return
		}

		respond(w, http.StatusCreated, result)
	}
}

func patchDepositsHandler(srv savingDeposits.DepositsService) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		defer r.Body.Close()

		var updateInput savingDeposits.DepositUpdateInput

		if err := json.NewDecoder(r.Body).Decode(&updateInput.Data); err != nil {
			respond(w, http.StatusBadRequest, err.Error())
			return
		}

		updateInput.Id = vars["id"]

		result, err := srv.Update(updateInput)
		if err != nil {
			badRequestError(err, w)
			return
		}

		respond(w, http.StatusOK, result)
	}
}

func deleteDepositsHandler(srv savingDeposits.DepositsService) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var deleteIn savingDeposits.DepositDeleteInput

		vars := mux.Vars(r)
		deleteIn.Id = vars["id"]
		_, err := srv.Delete(deleteIn)
		if err != nil {
			badRequestError(err, w)
			return
		}

		respond(w, http.StatusNoContent, nil)
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
