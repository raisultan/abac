CREATE TABLE IF NOT EXISTS users
(
    id SERIAL,
    email TEXT NOT NULL,
    firstName TEXT NOT NULL,
    lastName TEXT NOT NULL,
    createdAt TEXT NOT NULL,
    isAdmin BOOLEAN NOT NULL DEFAULT FALSE,
    isApproved BOOLEAN NOT NULL DEFAULT FALSE,

    CONSTRAINT users_pkey PRIMARY KEY (id)
);
