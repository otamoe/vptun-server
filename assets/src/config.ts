
let apiURL = ''
if (process.env.NODE_ENV === 'development') {
    apiURL = 'https://localhost:8443'
}


export const API_URL = apiURL
export const TITLE = "VPTun Server"
export const DESCRIPTION = "VPTun Server"
