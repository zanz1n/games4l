-- name: GetQuestionById :one
SELECT * FROM "question" WHERE "id" = $1;

-- name: CreateQuestion :one
INSERT INTO
    "question" (
        "question",
        "answer1",
        "answer2",
        "answer3",
        "answer4",
        "correct_answer",
        "type",
        "style",
        "difficulty"
    )
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING *;

-- name: GetMany :many
SELECT * FROM "question" LIMIT $1;
