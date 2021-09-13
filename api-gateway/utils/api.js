const axios = require('axios')

const { TIMEOUT } = process.env

module.exports = (baseURL) => {
    return axios.create({
        baseURL: baseURL,
        timeout: parseInt(TIMEOUT),
        headers: {
            'Content-type': 'application/json'
        }
    })
}