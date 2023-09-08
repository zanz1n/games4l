package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/games4l/internal/question"
	"github.com/goccy/go-json"
)

const (
	ModeOldToNew = "old-to-new"
	ModeNewToOld = "new-to-old"

	IndexTypeArray     = "array"
	IndexTypeStringMap = "map"
)

func ConvertQuestionOldToNew(inBuf []byte, indexType string) ([]byte, error) {
	var (
		err    error
		outBuf []byte
	)

	if indexType == IndexTypeArray {
		arr := []question.QuestionAlternativeFmt{}
		if err = json.Unmarshal(inBuf, &arr); err != nil {
			return nil, errors.New("failed to unmarshal input file: " + err.Error())
		}

		newArr := make([]*question.Question, len(arr))

		for i, question := range arr {
			if err = validate.Struct(&question); err != nil {
				return nil, fmt.Errorf("%d° input question is invalid: %s", i+1, err.Error())
			}

			if !question.IsValid() {
				return nil, fmt.Errorf("%d° input question is invalid", i+1)
			}

			// Implicit copy
			q := question

			newQuest := q.Parse()
			*newQuest.CorrectAnswer = *newQuest.CorrectAnswer - 1
			newArr[i] = newQuest
		}

		outBuf, err = json.Marshal(newArr)
	} else {
		mp := map[string]question.QuestionAlternativeFmt{}
		if err = json.Unmarshal(inBuf, &mp); err != nil {
			return nil, errors.New("failed to unmarshal input file: " + err.Error())
		}

		newArr := make([]*question.Question, len(mp))

		i := 0
		for key, question := range mp {
			if err = validate.Struct(&question); err != nil {
				return nil, fmt.Errorf("'%s' input question is invalid: %s", key, err.Error())
			}

			if !question.IsValid() {
				return nil, fmt.Errorf("'%s' input question is invalid", key)
			}

			// Implicit copy
			q := question

			newQuest := q.Parse()
			*newQuest.CorrectAnswer = *newQuest.CorrectAnswer - 1
			newArr[i] = newQuest
			i++
		}

		outBuf, err = json.Marshal(newArr)
	}

	if err != nil {
		return nil, errors.New("failed to marshal output data: " + err.Error())
	}

	return outBuf, nil
}

func HandleQuestionConvert(args map[string]string) error {
	mode, ok := args["mode"]
	if !ok {
		mode = ModeOldToNew
	} else if mode != ModeNewToOld && mode != ModeOldToNew {
		return errors.New("argument mode must be old-to-new or new-to-old")
	}

	indexType, ok := args["index-type"]
	if !ok {
		indexType = IndexTypeArray
	} else if indexType != IndexTypeArray && indexType != IndexTypeStringMap {
		return errors.New("argument index-type must be array or map")
	}

	inFile, ok := args["input"]
	if !ok {
		return errors.New("input argument must be provided")
	}
	outFile, ok := args["output"]
	if !ok {
		return errors.New("output argument must be provided")
	}

	inBuf, err := os.ReadFile(inFile)
	if err != nil {
		return errors.New("failed to read '" + inFile + "': " + err.Error())
	}

	var outBuf []byte
	if mode == ModeNewToOld {
		return errors.New("Amogus")
	} else {
		outBuf, err = ConvertQuestionOldToNew(inBuf, indexType)
	}

	if err = os.WriteFile(outFile, outBuf, 0644); err != nil {
		return errors.New("failed to write to file '" + outFile + "': " + err.Error())
	}

	return nil
}
