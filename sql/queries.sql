-- name: GetQuestionById :one
SELECT * FROM "question" WHERE "id" = $1;

-- name: CreateQuestion :one
INSERT INTO
    "question" (
        "question",
        "answer_1",
        "answer_2",
        "answer_3",
        "answer_4",
        "correct_answer",
        "type",
        "style",
        "difficulty"
    )
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING *;

-- name: GetManyQuestions :many
SELECT * FROM "question" ORDER BY "id" LIMIT $1;

-- name: UpdateQuestionById :one
UPDATE "question"
SET
    "updated_at" = CURRENT_TIMESTAMP,
    "question" = $1,
    "answer_1" = $2,
    "answer_2" = $3,
    "answer_3" = $4,
    "answer_4" = $5,
    "correct_answer" = $6,
    "type" = $7
WHERE "id" = $8 RETURNING *;

-- name: DeleteQuestionById :exec
DELETE FROM "question" WHERE "id" = $1;
