CREATE TABLE tormonger_data (
    link_hash TEXT NOT NULL,
    link TEXT NOT NULL,
    PRIMARY KEY(link_hash),
    UNIQUE(link)
);

CREATE TABLE html_data (
    id serial PRIMARY KEY,
    link_hash TEXT NOT NULL,
    html_data TEXT,
    CONSTRAINT fk_tormonger_data
       FOREIGN KEY(link_hash)
           REFERENCES tormonger_data(link_hash)
           ON DELETE CASCADE ON UPDATE cascade
);

CREATE TABLE tormonger_data_sub_directories (
    tormonger_data_id TEXT NOT NULL,
    html_data_id TEXT NOT NULL,
    PRIMARY KEY (tormonger_data_id, html_data_id),
    subdirectory_path TEXT,
    CONSTRAINT fk_tormonger_data
        FOREIGN KEY(tormonger_data_id)
            REFERENCES tormonger_data(link_hash)
            ON DELETE CASCADE ON UPDATE cascade,
    CONSTRAINT html_data
        FOREIGN KEY(html_data_id)
            REFERENCES html_data(id)
            ON DELETE CASCADE ON UPDATE cascade
);

CREATE TABLE logs (
    error_id INT generated always as identity,
    error_message VARCHAR(250),
    notes VARCHAR(250),
    PRIMARY KEY(error_id)
);