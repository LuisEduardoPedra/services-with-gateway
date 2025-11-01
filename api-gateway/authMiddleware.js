const jwt = require('jsonwebtoken');
const JWT_SECRET = process.env.JWT_SECRET;

const authMiddleware = (req, res, next) => {
  const authHeader = req.headers.authorization;
  if (!authHeader) {
    return res.status(401).json({ error: 'Token de autorização não fornecido' });
  }

  const parts = authHeader.split(' ');
  if (parts.length !== 2 || parts[0] !== 'Bearer') {
    return res.status(401).json({ error: 'Formato do token inválido' });
  }

  const token = parts[1];

  try {

    const claims = jwt.verify(token, JWT_SECRET);

    req.user = claims; 
    next();
  } catch (err) {
    return res.status(401).json({ error: 'Token inválido ou expirado' });
  }
};

const permissionMiddleware = (requiredPermission) => {
  return (req, res, next) => {

    if (!req.user || !req.user.roles) {
      return res.status(403).json({ error: 'Claims do usuário não encontrados' });
    }

    const roles = req.user.roles;

    if (roles.includes(requiredPermission)) {
      next(); 
    } else {
      return res.status(403).json({ error: 'Acesso negado: permissão necessária ausente' });
    }
  };
};

module.exports = {
  authMiddleware,
  permissionMiddleware
};