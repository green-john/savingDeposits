import {$http} from "./http";
import $auth from "./auth";

export default {
    loadAllDeposits() {
        let finalUrl = '/deposits';
        // if (filters.floorAreaMeters) {
        //     finalUrl += `floorAreaMeters=${filters.floorAreaMeters}&`;
        // }
        //
        // if (filters.pricePerMonthUSD) {
        //     finalUrl += `pricePerMonthUSD=${filters.pricePerMonthUSD}&`;
        // }
        //
        // if (filters.roomCount) {
        //     finalUrl += `roomCount=${filters.roomCount}`;
        // }

        return $http.get(
            finalUrl, {
                headers: {Authorization: $auth.getToken()}
            }
        ).then(response => {
            console.log('loadAll');
            console.log(response);
            return response.data;
        }).catch(err => {
            throw err;
        })
    },

    createDeposit(depositData) {
        console.log('create');
        console.log(depositData);
        // depositData.startDate = dateToStr(depositData.startDate);
        // depositData.endDate = dateToStr(depositData.endDate);

        return $http.post('/deposits', depositData, {
            headers: {Authorization: $auth.getToken()}
        }).then(response => {
            return response.data;
        }).catch(err => {
            throw err;
        })
    },

    deleteDeposit(depositId) {
        return $http.delete('/deposits/' + depositId, {
            headers: {Authorization: $auth.getToken()}
        }).then(res => {
            return res.data;
        }).catch(err => {
            throw err;
        });
    },

    updateDeposit(id, newData) {
        console.log('update');
        console.log(newData);

        return $http.patch('/deposits/' + newData.id, newData, {
            headers: {Authorization: $auth.getToken()}
        }).then(response => {
            return response.data;
        }).catch(err => {
            throw err;
        })
    }
}