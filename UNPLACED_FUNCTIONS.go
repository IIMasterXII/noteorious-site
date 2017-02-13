package main
import (
	"strings"
	"errors"
	"net/http"
	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"strconv"
	"github.com/Esseh/retrievable"
	"google.golang.org/appengine/datastore"
)


type Context struct {
	req *http.Request
	res http.ResponseWriter
	user *User
	userException error
	context.Context
}

func (ctx Context)AssertLoggedInFailed() bool {
	if ctx.userException != nil {
		path := strings.Replace(ctx.req.URL.Path[1:], "%2f", "/", -1)
		http.Redirect(ctx.res, ctx.req, PATH_AUTH_Login+"?redirect="+path, http.StatusSeeOther)
		return true
	}
	return false
}

func NewContext(res http.ResponseWriter, req *http.Request) Context{
	user, err := GetUserFromSession(req)
	ctx := Context { 
		req: req,
		res: res,
		user: user,
		userException: err,
	}
	ctx.Context = appengine.NewContext(req)
	return ctx
}

func GetExistingNote(id string, ctx context.Context)(*Note,*Content,error){
	RetrievedNote := &Note{}
	RetrievedContent := &Content{}
	NoteKey, err := strconv.ParseInt(id, 10, 64)
	if err != nil { return RetrievedNote,RetrievedContent,err }
	err = retrievable.GetEntity(ctx, NoteKey, RetrievedNote)
	if err != nil { return RetrievedNote,RetrievedContent,err }
	err = retrievable.GetEntity(ctx, RetrievedNote.ContentID, RetrievedContent)
	return RetrievedNote,RetrievedContent,err
}

func VerifyNotePermission(res http.ResponseWriter, req *http.Request, user *User, note *Note) bool {
	redirect := strconv.FormatInt(note.OwnerID, 10)
	if note.OwnerID != int64(user.IntID) && note.Protected {
		http.Redirect(res, req, "/view/"+redirect, http.StatusSeeOther)
		return false
	}
	return true
}

func CreateNewNote(ctx context.Context,NewContent Content,NewNote Note) (*datastore.Key,*datastore.Key,error) {
	contentKey, err := retrievable.PlaceEntity(ctx, int64(0), &NewContent)
	if err != nil { return contentKey,&datastore.Key{},err }
	NewNote.ContentID = contentKey.IntID()
	noteKey, err := retrievable.PlaceEntity(ctx, int64(0), &NewNote)
	return contentKey,noteKey,err
}

func UpdateNoteContent(ctx context.Context, res http.ResponseWriter,req *http.Request,user *User,id string,UpdatedContent Content, UpdatedNote Note) error {
	Note := Note{}
	noteID, err := strconv.ParseInt(id, 10, 64); if err != nil { return err }
	err = retrievable.GetEntity(ctx, noteID, &Note); if err != nil { return err }
	validated := VerifyNotePermission(res, req, user, &Note); if !validated { return errors.New("Permission Error: Not Allowed") }
	if Note.OwnerID == int64(user.IntID) { Note.Protected = UpdatedNote.Protected }	
	_, err = retrievable.PlaceEntity(ctx, noteID, &Note); if err != nil { return err }
	_, err = retrievable.PlaceEntity(ctx, Note.ContentID, &UpdatedContent); return err
}