CREATE TABLE tormonger_data (
    id uuid DEFAULT gen_random_uuid() PRIMARY KEY NOT NULL,
    link_hash TEXT NOT NULL,
    link TEXT NOT NULL,
    UNIQUE(link)
);

-- one tormonger_data row to many tormonger_data_sub_directories
CREATE TABLE tormonger_data_sub_directories (
    id uuid DEFAULT gen_random_uuid() PRIMARY KEY NOT NULL,
    tormonger_data_id uuid NOT NULL,
    subdirectory_path TEXT,
    CONSTRAINT fk_tormonger_data
        FOREIGN KEY(tormonger_data_id)
            REFERENCES tormonger_data(id)
            ON DELETE CASCADE ON UPDATE cascade
);

-- one tormonger_data row to many html_data rows due to sub dirs
CREATE TABLE html_data (
    id uuid DEFAULT gen_random_uuid() PRIMARY KEY NOT NULL,
    tormonger_data_id uuid NOT NULL,
    tormonger_data_sub_directories_id uuid,
    html_data TEXT,
    CONSTRAINT fk_tormonger_data
       FOREIGN KEY(tormonger_data_id)
           REFERENCES tormonger_data(id)
           ON DELETE CASCADE ON UPDATE cascade,
    CONSTRAINT fk_tormonger_data_sub_directories
        FOREIGN KEY(tormonger_data_sub_directories_id)
            REFERENCES tormonger_data_sub_directories(id)
            ON DELETE CASCADE ON UPDATE cascade
);

CREATE TABLE logs (
    id uuid DEFAULT gen_random_uuid() PRIMARY KEY NOT NULL,
    log_message VARCHAR(250),
    log_type VARCHAR(250)
);