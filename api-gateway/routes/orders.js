const express = require('express')
const router = express.Router()

const apiAdapter = require('../utils/api')
const {URL_ORDER_SERVICE} = process.env
const api = apiAdapter(URL_ORDER_SERVICE)

router.post('/', async (req, res) => {
    try {
        const book = await api.post('/orders', req.body)
        return res.json(book.data)
    } catch (error) {
        if (error.code === 'ECONNREFUSED') {
            return res.status(500).json({
                status: 'error',
                message: 'Service unavailable!'
            })
        }

        const { status, data } = error.response
        return res.status(status).json(data)
    }
})

module.exports = router