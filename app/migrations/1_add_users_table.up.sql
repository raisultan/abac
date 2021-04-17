CREATE TABLE IF NOT EXISTS users
(
    id SERIAL,
    email TEXT NOT NULL,
    firstName TEXT NOT NULL,
    lastName TEXT NOT NULL,
    createdAt TEXT DEFAULT NULL,
    isAdmin BOOLEAN DEFAULT FALSE,
    isApproved BOOLEAN DEFAULT FALSE,

    CONSTRAINT users_pkey PRIMARY KEY (id)
);
