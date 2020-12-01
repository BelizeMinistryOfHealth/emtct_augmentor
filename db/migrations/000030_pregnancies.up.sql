DROP TABLE pregnancies;
CREATE TABLE pregnancies(
    pregnancy_id int PRIMARY KEY ,
    patient_id int NOT NULL,
    lmp DATE,
    edd DATE,
    end_time DATE
);
