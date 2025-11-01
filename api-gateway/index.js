// Load environment variables from .env file
require('dotenv').config();

const express = require('express');
const cors = require('cors');
const { createProxyMiddleware } = require('http-proxy-middleware');
const { authMiddleware, permissionMiddleware } = require('./authMiddleware');

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

//const authServiceTarget = 'http://localhost:8081';
//const analysisServiceTarget = 'http://localhost:8082';
//const converterServiceTarget = 'http://localhost:8083';

const authServiceTarget = 'http://auth-service:8081';
const analysisServiceTarget = 'http://analysis-service:8082';
const converterServiceTarget = 'http://converter-service:8083';

app.use('/api/v1/login', createProxyMiddleware({
  target: authServiceTarget,
  changeOrigin: true,

  pathRewrite: {
    '^/': '/api/v1/login',
  },
  
  onProxyReq: (proxyReq, req, res) => {
    console.log(`[Gateway] Proxying to Auth Service: ${req.method} ${req.path}`);
  }
}));

app.use(
  '/api/v1/analyze/icms',
  authMiddleware, 
  permissionMiddleware('analise-icms'), 
  createProxyMiddleware({ 
    target: analysisServiceTarget,
    changeOrigin: true,
    pathRewrite: {
    '^/': '/api/v1/analyze/icms',
    },
    onProxyReq: (proxyReq, req, res) => {
      console.log(`[Gateway] Proxying to Analysis Service: ${req.method} ${req.path}`);
    }
  })
);

app.use(
  '/api/v1/analyze/ipi-st',
  authMiddleware, 
  permissionMiddleware('analise-ipi-st'), 
  createProxyMiddleware({ 
    target: analysisServiceTarget,
    changeOrigin: true,
    pathRewrite: {
    '^/': '/api/v1/analyze/ipi-st',
    },
    onProxyReq: (proxyReq, req, res) => {
      console.log(`[Gateway] Proxying to Analysis Service: ${req.method} ${req.path}`);
    }
  })
);


app.use(
  '/api/v1/convert/francesinha',
  authMiddleware,
  permissionMiddleware('converter-francesinha'),
  createProxyMiddleware({
    target: converterServiceTarget,
    changeOrigin: true,
    pathRewrite: {
    '^/': '/api/v1/convert/francesinha',
    },
    onProxyReq: (proxyReq, req, res) => {
      console.log(`[Gateway] Proxying to Converter Service: ${req.method} ${req.path}`);
    }
  })
);

app.use(
  '/api/v1/convert/receitas-acisa',
  authMiddleware,
  permissionMiddleware('converter-receitas-acisa'),
  createProxyMiddleware({
    target: converterServiceTarget,
    changeOrigin: true,
    pathRewrite: {
    '^/': '/api/v1/convert/receitas-acisa',
    },
    onProxyReq: (proxyReq, req, res) => {
      console.log(`[Gateway] Proxying to Converter Service: ${req.method} ${req.path}`);
    }
  })
);

app.use(
  '/api/v1/convert/atolini-pagamentos',
  authMiddleware,
  permissionMiddleware('converter-atolini-pagamentos'),
  createProxyMiddleware({
    target: converterServiceTarget,
    changeOrigin: true,
    pathRewrite: {
    '^/': '/api/v1/convert/atolini-pagamentos',
    },
    onProxyReq: (proxyReq, req, res) => {
      console.log(`[Gateway] Proxying to Converter Service: ${req.method} ${req.path}`);
    }
  })
);

app.use(
  '/api/v1/convert/atolini-recebimentos',
  authMiddleware,
  permissionMiddleware('converter-atolini-recebimentos'),
  createProxyMiddleware({
    target: converterServiceTarget,
    changeOrigin: true,
    pathRewrite: {
    '^/': '/api/v1/convert/atolini-recebimentos',
    },
    onProxyReq: (proxyReq, req, res) => {
      console.log(`[Gateway] Proxying to Converter Service: ${req.method} ${req.path}`);
    }
  })
);

app.get('/api/v1/health', (req, res) => {
  res.json({ status: 'UP', message: 'Gateway is running!' });
});

app.listen(port, () => {
  console.log(`🚀 API Gateway server listening on http://localhost:${port}`);
});