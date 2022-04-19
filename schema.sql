create table if not exists products(
    id SERIAL PRIMARY KEY,
    guid character varying(255) UNIQUE NOT NULL,
    name character varying(255) UNIQUE NOT NULL,
    price REAL NOT NULL,
    description character varying,
    createdAt character varying NOT NULL
)
