package main

import (
	"html/template"
	"math/rand"
	"strconv"
	"strings"

	"github.com/Esseh/notorious-dev/CONTEXT"
	"github.com/Esseh/notorious-dev/COOKIE"
	"github.com/Esseh/notorious-dev/CORE"
	"github.com/Esseh/notorious-dev/NOTES"
	"github.com/Esseh/notorious-dev/USERS"
	"github.com/Esseh/retrievable"
	humanize "github.com/dustin/go-humanize" // russross markdown parser
	appcontext "golang.org/x/net/context"
)

func init() {
	// Tie functions into template here with ... "functionName":theFunction,
	funcMap := template.FuncMap{
		"getAvatarURL":  CORE.GetAvatarURL,
		"getUser":       GetUserFromID,
		"humanize":      humanize.Time,
		"humanizeSize":  humanize.Bytes,
		"monthfromtime": CORE.MonthFromTime,
		"yearfromtime":  CORE.YearFromTime,
		"dayfromtime":   CORE.DayFromTime,
		"findsvg":       CORE.FindSVG,
		"findtemplate":  CORE.FindTemplate,
		"inc":           CORE.Inc,
		"addCtx":        CORE.AddCtx,
		"getDate":       CORE.GetDate,
		"toInt":         CORE.ToInt,
		"getMod":        GetMod,
		"canEditNote":	 NOTES.CanEditNote,
		"canViewNote":	 NOTES.CanViewNote,
		"findCollabs":   FindCollaborators,
		"canEdit":       NOTES.CanEditNote,
		"canView":       NOTES.CanViewNote,
		// "isOwner":       isOwner,
		"parse": CORE.EscapeString,
	} // Load up all templates.
	CORE.TPL = template.New("").Funcs(funcMap)
	CORE.TPL = template.Must(CORE.TPL.ParseGlob("templates/*"))
}

func GetMod(a int64) int64 {
	rand.Seed(a)
	return int64(rand.Uint32()) % 10
}

func FindCollaborators(c string) []int64 {
	dupCheck := make(map[int]bool)
	temp := strings.Split(c, ":")
	var collabs []int64
	for _, x := range temp {
		collab, err := strconv.Atoi(x)
		if (err == nil) && (!dupCheck[collab]) {
			collabs = append(collabs, int64(collab))
			dupCheck[collab] = true
		}
	}
	return collabs
}

func GetUserFromID(ctx appcontext.Context, id int64) (*USERS.User, error) {
	owner := &USERS.User{}
	err := retrievable.GetEntity(ctx, id, owner)
	if err != nil {
		return &USERS.User{}, err
	} else {
		return owner, nil
	}
}

// Constructs the header.
// As the header gets more complex(such as capturing the current path)
// the need for such a helper function increases.
func MakeHeader(ctx CONTEXT.Context) *CONTEXT.HeaderData {
	oldCookie, err := COOKIE.GetValue(ctx.Req, "session")
	if err == nil {
		COOKIE.Make(ctx.Res, "session", oldCookie)
	}
	redirectURL := ctx.Req.URL.Path[1:]
	if redirectURL == "login" || redirectURL == "register" || redirectURL == "elevatedlogin" {
		redirectURL = ctx.Req.URL.Query().Get("redirect")
	}
	return &CONTEXT.HeaderData{
		ctx, ctx.User, redirectURL,
	}
}
