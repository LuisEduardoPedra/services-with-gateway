// Load environment variables from .env file
// This replaces the 'loadEnv' function from your Go main.go
require('dotenv').config();

const express = require('express');
const cors = require('cors');

const app = express();

const port = process.env.PORT || 8080; 

const allowedOrigins = (process.env.ALLOWED_ORIGINS || '').split(',');

const corsOptions = {
  origin: function (origin, callback) {
    if (!origin || allowedOrigins.indexOf(origin) !== -1 || allowedOrigins.indexOf('*') !== -1) {
      callback(null, true);
    } else {
      callback(new Error('Not allowed by CORS'));
    }
  }
};

app.use(cors(corsOptions));

app.get('/api/v1/health', (req, res) => {
  res.json({ status: 'UP', message: 'Gateway is running!' });
});

app.listen(port, () => {
  console.log(`🚀 API Gateway server listening on http://localhost:${port}`);
});