const express = require('express')
const router = express.Router()

const apiAdapter = require('../utils/api.js')
const { URL_BOOK_SERVICE } = process.env
const api = apiAdapter(URL_BOOK_SERVICE)

router.get('/', async (req, res) => {
    try {
        const books = await api.get('/books')
        return res.json(books.data)
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

router.get('/:id', async (req, res) => {
    try {
        const id = req.params.id
        const book = await api.get(`/books/${id}`)
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

router.post('/', async (req, res) => {
    try {
        const book = await api.post('/books', req.body)
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

router.put('/:id', async (req, res) => {
    try {
        const id = req.params.id
        const book = await api.put(`/books/${id}`, req.body)
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

router.delete('/:id', async (req, res) => {
    try {
        const id = req.params.id
        const book = await api.delete(`/books/${id}`)
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