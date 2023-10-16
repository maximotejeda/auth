CREATE TABLE IF NOT EXISTS users(
       id INTEGER PRIMARY KEY,
       username TEXT UNIQUE NOT NULL,
       password TEXT NOT NULL,
       email TEXT UNIQUE,
       rol TEXT,
       created_at INTEGER,
       edited_at INTEGER
);

CREATE TABLE IF NOT EXISTS login_info(
       id INTEGER PRIMARY KEY,
       user_id INTEGER,
       token TEXT,
       login_list TEXT,
       ip_list TEXT,
       user_agen_list TEXT,
       failed_login_list TEXT,
       mac_address_list TEXT,
       created_at INTEGER,
       edited_at INTEGER,	
       FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE    
); 

