-- Basic table may add EditedAt to keep track of last time updated
CREATE TABLE IF NOT EXISTS Tasks (
    TaskID INTEGER PRIMARY KEY AUTOINCREMENT,
    TaskName TEXT NOT NULL,     
    TaskDescription TEXT NOT NULL,     
    IsCompleted BOOLEAN NOT NULL DEFAULT 0,
    CreatedAt DATETIME DEFAULT CURRENT_TIMESTAMP
);
