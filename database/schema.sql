CREATE TABLE galaxies(
    id                  INT AUTO_INCREMENT NOT NULL,
    name                VARCHAR(200) NOT NULL,
    ugc_number          VARCHAR(200) NOT NULL,
    created_at          TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    PRIMARY KEY (`id`),
    UNIQUE (ugc_number)
);

CREATE TABLE stars(
    id                  INT AUTO_INCREMENT NOT NULL,
    galaxy_id           INT NOT NULL,
    name                VARCHAR(200) NOT NULL,
    gaia_catalogue_id   VARCHAR(200) NOT NULL,
    created_at          TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    PRIMARY KEY (`id`)
);
