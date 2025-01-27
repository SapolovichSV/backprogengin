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
INSERT INTO drinks(name, tags)
VALUES ("Bad Touch","Sour,Classy,Vintage"),
    ("Beer", "Bubbly,Classic,Vintage"),
    ("Bleeding Jane", "Spicy,Classic,Sobering"),
    ("Bloom Light", "Spicy,Promo,Bland"),
    ("Blue Fairy", "Sweet,Girly,Soft"),
    ("Brandtini", "Sweet,Classy,Happy"),
    ("Cobalt Velvet", "Bubbly,Classy,Burning"),
    ("Crevice Spike", "Sour,Manly,Sobering"),
    ("Flaming Moai", "Sour,Classy"),
    ("Fluffy Dream", "Sour,Girly,Soft"),
    ("Fringe Weaver", "Bubbly,Classy,Strong"),
    ("Frothy Water", "Bubbly,Classic,Bland"),
    ("Grizzly Temple", "Bitter,Promo,Bland"),
    ("Gut Punch", "Bitter,Manly,Strong"),
    ("Marsblast", "Spicy,Manly,Strong"),
    ("Mercuryblast", "Sour,Classy,Burning"),
    ("Moonblast", "Sweet,Girly,Happy"),
    ("Piano Man", "Sour,Promo,Strong"),
    ("Piano Woman", "Sweet,Promo,Happy"),
    ("Pile Driver", "Bitter,Manly,Burning"),
    ("Sparkle Star", "Sweet,Girly,Happy"),
    ("Sugar Rush", "Sweet,Girly,Happy"),
    ("Sunshine Cloud", "Bitter,Girly,Soft"),
    ("Suplex", "Bitter,Manly,Burning"),
    ("Zen Star", "Sour,Promo,Bland");
