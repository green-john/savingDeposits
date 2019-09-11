import {$http} from "./http";
import $auth from "./auth";

export default {
    loadAllDeposits(filters) {
        let finalUrl = '/deposits?';

        if (filters.minAmount) {
            finalUrl += `minAmount=${filters.minAmount}&`;
        }

        if (filters.maxAmount) {
            finalUrl += `maxAmount=${filters.maxAmount}&`;
        }

        if (filters.bankName) {
            finalUrl += `bankName=${filters.bankName}&`;
        }

        if (filters.startDate) {
            finalUrl += `startDate=${filters.startDate}&`;
        }

        if (filters.endDate) {
            finalUrl += `endDate=${filters.endDate}&`;
        }

        return $http.get(
            finalUrl, {
                headers: {Authorization: $auth.getToken()}
            }
        ).then(response => {
            return response.data;
        });
    },

    createDeposit(depositData) {

        return $http.post('/deposits', depositData, {
            headers: {Authorization: $auth.getToken()}
        }).then(response => {
            return response.data;
        });
    },

    deleteDeposit(depositId) {
        return $http.delete('/deposits/' + depositId, {
            headers: {Authorization: $auth.getToken()}
        }).then(res => {
            return res.data;
        });
    },

    updateDeposit(id, newData) {
        return $http.patch('/deposits/' + newData.id, newData, {
            headers: {Authorization: $auth.getToken()}
        }).then(response => {
            return response.data;
        });
    },

    getReport(startDate, endDate) {
        let finalUrl = "/report?";

        if (startDate) {
            finalUrl += `startDate=${startDate}&`;
        }

        if (endDate) {
            finalUrl += `endDate=${endDate}&`;
        }

        return $http.get(finalUrl, {
            headers: {Authorization: $auth.getToken()}
        }).then(res => {
            return res.data;
        });
    }
}