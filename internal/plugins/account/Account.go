package account

import (
	"fmt"
	"myserver/internal/connection"
	"myserver/internal/database"
	"net/http"

	"github.com/google/uuid"
)

var AccountRoutePath = []connection.RoutePathHandler{}

func createTransaction(rh *connection.RouteHandler) error {
	fmt.Println("Add transaction")
	req := database.CreateTransaction{}

	if err := rh.GetJSON(&req); err != nil {
		rh.ResponseError(http.StatusBadRequest, "Invalid body")
		return nil
	}

	accountCon := rh.DbHandler.GetAccountConnection()

	transactionId := accountCon.AddTransaction(req)

	rh.ResponseJSON(http.StatusOK, map[string]any{
		"id": transactionId,
	})

	return nil
}

func CreateNewNotes(rh *connection.RouteHandler) error {
	fmt.Println("CreateNewNotes")
	req := database.CreateNotes{}

	if err := rh.GetJSON(&req); err != nil {
		rh.ResponseError(http.StatusBadRequest, "Invalid body")
		return nil
	}

	noteCon := rh.DbHandler.GetNotesConnection()

	notesId := noteCon.CreateNotes(req)

	rh.ResponseJSON(http.StatusOK, map[string]any{
		"id": notesId,
	})

	return nil
}

func CreateNewNoteGroup(rh *connection.RouteHandler) error {
	fmt.Println("CreateNewNoteGroup")
	req := database.CreateGroupNote{}

	if err := rh.GetJSON(&req); err != nil {
		rh.ResponseError(http.StatusBadRequest, "Invalid body")
		return nil
	}

	noteCon := rh.DbHandler.GetNotesConnection()

	notesId := noteCon.CreateGroupNote(req)

	rh.ResponseJSON(http.StatusOK, map[string]any{
		"id": notesId,
	})

	return nil
}

func UpdateNotes(rh *connection.RouteHandler) error {
	id, err := uuid.Parse(rh.GetPathValue("id"))
	if err != nil {
		rh.ResponseError(http.StatusBadRequest, "Invalid UUID")
		return err
	}
	req := database.UpdateNotes{}

	if err := rh.GetJSON(&req); err != nil {
		rh.ResponseError(http.StatusBadRequest, "Invalid body")
		return nil
	}

	fmt.Println(req.NoteGroup)

	noteCon := rh.DbHandler.GetNotesConnection()
	noteCon.UpdateNotes(id, req.Text, req.Title, req.NoteGroup)

	rh.ResponseJSON(http.StatusOK, map[string]any{
		"id": id,
	})

	return nil
}

func GetMyNotes(rh *connection.RouteHandler) error {
	fmt.Println("Get user")
	userId, err := uuid.Parse(rh.GetPathValue("userId"))
	if err != nil {
		rh.ResponseError(http.StatusBadRequest, "Invalid UUID")
		return err
	}
	noteCon := rh.DbHandler.GetNotesConnection()
	notes := noteCon.GetAllNotes(userId)

	rh.ResponseJSON(http.StatusOK, notes)

	return nil
}

func GetMyNotesById(rh *connection.RouteHandler) error {
	fmt.Println("Get notes by id")
	id, err := uuid.Parse(rh.GetPathValue("id"))
	if err != nil {
		rh.ResponseError(http.StatusBadRequest, "Invalid UUID")
		return err
	}
	noteCon := rh.DbHandler.GetNotesConnection()
	notes := noteCon.GetNotesById(id)

	rh.ResponseJSON(http.StatusOK, notes)

	return nil
}
