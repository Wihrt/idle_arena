db = db.getSiblingDB('admin')
db.createUser(
    {
        user: "arenaManager",
        pwd: "u2MQrvfHbTLxqjYf",
        roles: [ { role: "readWrite", db: "arena" } ]
    }
)