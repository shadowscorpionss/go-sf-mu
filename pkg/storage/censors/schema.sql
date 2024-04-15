
DROP TABLE IF EXISTS black_list;

CREATE TABLE IF NOT EXISTS black_list (
     id SERIAL PRIMARY KEY,
     ban_word TEXT
);


INSERT INTO black_list (ban_word) VALUES ('qwerty');

INSERT INTO black_list (ban_word) VALUES ('йцукен');

INSERT INTO black_list (ban_word) VALUES ('zxvbnm');
