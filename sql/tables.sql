/* Create your tables here. You can get the script from the other database
and make adjustments to use in psql. */

CREATE TABLE TABLE_1 (
    ID           INTEGER GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
    CADASTRO     TIMESTAMP WITHOUT TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    AGENDA_DATA  DATE NOT NULL,
    AGENDA_HORA  TIME WITHOUT TIME ZONE NOT NULL,
    MSG          TEXT NOT NULL
);
