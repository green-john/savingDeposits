import {$http} from "./http";
import $auth from "./auth";

function formatPercentForDisplay(deposits) {
    for (let deposit of deposits) {
        deposit.tax *= 100;
        deposit.yearlyInterest *= 100;
    }

    return deposits;
}

function formatForServer(deposit) {
    if (deposit.tax) {
        deposit.tax /= 100;
    }

    if (deposit.yearlyInterest) {
        deposit.yearlyInterest /= 100;
    }

    return deposit;
}

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
            return formatPercentForDisplay(response.data);
        });
    },

    createDeposit(depositData) {
        depositData = formatForServer(depositData);
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
        newData = formatForServer(newData);
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
            return formatPercentForDisplay(res.data);
        });
    }
}