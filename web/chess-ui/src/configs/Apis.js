import axios from 'axios'

const BASE_DOMAIN = 'localhost:8080'

const BASE_API_URL = `http://${BASE_DOMAIN}/api/v1`

export const BASE_WS_URL = `ws://${BASE_DOMAIN}/ws`

export const endpoint = {
    'auth': '/auth',
}

export default axios.create({
    baseURL: BASE_API_URL
})