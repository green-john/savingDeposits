import {$http} from "./http";
import $auth from "./auth";

export default {
    createClientAccount(username, password) {
        return $http.post('/newClient', {username, password}).then(
            res => {
                return res.data || [];
            }
        ).catch(err => {
            throw err;
        });
    },

    createUser(username, password, role) {
        return $http.post("/users", {username, password, role}, {
            headers: {Authorization: $auth.getToken()}
        }).then(
            res => {
                return res.data || [];
            }
        ).catch(err => {
            throw err;
        });
    },

    updateUser(id, username, password, role) {
        const userUrl = `/users/${id}`;
        return $http.patch(userUrl, {username, password, role}, {
            headers: {Authorization: $auth.getToken()}
        }).then(
            res => {
                return res.data || [];
            }
        ).catch(err => {
            throw err;
        });
    },

    getUser(id) {
        const userUrl = `/users/${id}`;
        return $http.get(userUrl).then(
            res => {
                return res.data || [];
            }
        ).catch(err => {
            throw err;
        });
    },

    deleteUser(id) {
        const userUrl = `/users/${id}`;
        return $http.delete(userUrl, {
            headers: {Authorization: $auth.getToken()}
        }).then(
            res => {
                return res.data || [];
            }
        ).catch(err => {
            throw err;
        });
    },

    getAllUsers() {
        return $http.get('/users', {
            headers: {Authorization: $auth.getToken()}
        }).then(res => {
            return res.data || [];
        }).catch(err => {
            throw err;
        })
    }
}