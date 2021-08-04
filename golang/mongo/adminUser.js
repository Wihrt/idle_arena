db = db.getSiblingDB('admin')
db.createUser(
    {
        user: _getEnv('MONGODB_USER'),
        pwd: _getEnv('MONGODB_PASSWORD'),
        roles: [ { role: "readWrite", db: "arena" } ]
    }
)