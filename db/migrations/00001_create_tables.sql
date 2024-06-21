-- +goose up
CREATE TABLE IF NOT EXISTS page_views (
  page_id VARCHAR(255) NOT NULL PRIMARY KEY, -- 255 is the max length of a folder name
  views INT NOT NULL DEFAULT 0
);
CREATE TABLE IF NOT EXISTS page_views_ips ( -- Used to ignore multiple views from the same IP
  page_id VARCHAR(255) NOT NULL,
  ip VARCHAR(255) NOT NULL CHECK(ip <> ''),
  PRIMARY KEY (page_id, ip),
  FOREIGN KEY (page_id) REFERENCES page_views(page_id)
);

-- +goose down
DROP TABLE IF EXISTS page_views_ips;
DROP TABLE IF EXISTS page_views;
