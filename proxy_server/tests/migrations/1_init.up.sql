CREATE TABLE thumbnails (
                            id INTEGER PRIMARY KEY AUTOINCREMENT, -- Уникальный идентификатор
                            url TEXT UNIQUE NOT NULL,       -- Ссылка на видео (уникальная)
                            data BLOB                      -- Превью в виде байтов
);
CREATE INDEX IF NOT EXISTS idx_url ON thumbnails (url);