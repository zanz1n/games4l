CREATE TYPE "QuestionType" AS ENUM ('2Alt', '4Alt');

CREATE TYPE "QuestionStyle" AS ENUM (
    'image',
    'audio',
    'video',
    'text'
);

CREATE TABLE
    "question" (
        "id" SERIAL PRIMARY KEY,
        "created_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
        "updated_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
        "question" VARCHAR(200) NOT NULL,
        "answer_1" VARCHAR(64) NOT NULL,
        "answer_2" VARCHAR(64) NOT NULL,
        "answer_3" VARCHAR(64),
        "answer_4" VARCHAR(64),
        "correct_answer" SMALLINT NOT NULL,
        "type" "QuestionType" NOT NULL,
        "style" "QuestionStyle" NOT NULL,
        "difficulty" SMALLINT NOT NULL
    );
