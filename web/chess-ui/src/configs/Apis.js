import axios from 'axios'

const BASE_URL = 'http://localhost:8080/api/v1'

export const endpoint = {
    'auth': '/auth',
}

export default axios.create({
    baseURL: BASE_URL
})