-- Selecciona la base de datos
USE db_name;

-- Crear tabla 'extra'
CREATE TABLE IF NOT EXISTS extra (
    id INT NOT NULL,
    extra_col_1 TEXT,
    CONSTRAINT pk_extra PRIMARY KEY (id)
);

-- Crear tabla 'other'
CREATE TABLE IF NOT EXISTS other (
    id INT NOT NULL,
    other_col_1 TEXT,
    other_col_2 TEXT,
    CONSTRAINT pk_other PRIMARY KEY (id)
);

-- Crear tabla 'test' con claves foráneas
CREATE TABLE IF NOT EXISTS test (
    id INT NOT NULL,
    col_1 TEXT,
    col_2 TEXT,
    other_id INT,
    extra_id INT,
    CONSTRAINT pk_test PRIMARY KEY (id),
    CONSTRAINT fk_other FOREIGN KEY (other_id) REFERENCES other(id),
    CONSTRAINT fk_extra FOREIGN KEY (extra_id) REFERENCES extra(id)
);

-- Insertar datos en 'extra'
INSERT INTO extra (id, extra_col_1) VALUES
(1, 'Extra info 1'),
(2, 'Extra info 2');

-- Insertar datos en 'other'
INSERT INTO other (id, other_col_1, other_col_2) VALUES
(10, 'Other info A', 'Details A'),
(20, 'Other info B', 'Details B');

-- Insertar datos en 'test'
INSERT INTO test (id, col_1, col_2, other_id, extra_id) VALUES
(100, 'Test Row 1', 'Some text', 10, 1),
(200, 'Test Row 2', 'More text', 20, 2);

