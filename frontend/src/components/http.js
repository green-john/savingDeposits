import Axios from 'axios';

const baseUrl = process.env.VUE_APP_API_ADDR;
export let $http = Axios.create({baseURL: baseUrl});

export function handleError(err) {
    if (err.response) {
        console.log(err.response.status);
        alert(`[ERROR] ${err.response.data}`);
    } else if (err.request) {
        alert(`[ERROR] ${err.request}`);
    } else {
        alert(`[ERROR] ${err.message}`);
    }

    console.log(err.config);
}
