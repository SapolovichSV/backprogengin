CREATE TABLE drinks (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    tags TEXT
);
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(255) NOT NULL,
    password VARCHAR(255)
);
CREATE TABLE favs (
    user_id INT NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id)  ON DELETE CASCADE,
    drink_id INT NOT NULL,
   FOREIGN KEY (drink_id) REFERENCES drinks(id) ON DELETE CASCADE
);
