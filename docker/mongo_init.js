use athena;
db.createUser({
    user: "athena_srvc",
    pwd: "<REPLACE ME>",
    roles: [
        { role: "readWrite", db: "athena" }
    ]
});
