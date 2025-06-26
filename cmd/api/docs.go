package main

import (
	"net/http"
	"strings"
	"time"

	"github.com/anvidev/apiduck"
	"github.com/anvidev/project-time-tracker/internal/store/categories"
	"github.com/anvidev/project-time-tracker/internal/store/sessions"
	"github.com/anvidev/project-time-tracker/internal/store/time_entries"
	"github.com/anvidev/project-time-tracker/internal/store/users"
	"github.com/anvidev/project-time-tracker/internal/types"
)

func ptr[T any](v T) *T {
	p := new(T)
	*p = v
	return p
}

func initDocumentation(config Config) *apiduck.Documentation {
	docs := apiduck.New(
		"Tidsregistrering API",
		"Internt værktøj til at dokumentere og overskue den tid der er brugt på projekter",
		config.Server.Version,
		apiduck.WithContact(
			"Skancode A/S",
			"support@skancode.dk",
			"www.skancode.dk",
		),
	)

	docs.AddServer("http://localhost:9090", "Development Server")
	docs.AddServer("https://api.tid.skancode.dk", "Production Server")

	docs.AddSecurity(apiduck.BearerToken("(bearer-token-for-users)", "Brugere skal være logget ind for at få adgang til ressourcer med denne authentication"))

	authResource := docs.AddResource("Auth", "Authentication og authorization")

	authResource.Post("/v1/auth/register", "Opret bruger", "Opret en ny bruger").
		Body(
			apiduck.JSONBody(users.RegisterUserInput{}).Example(users.RegisterUserInput{
				Name:     "John Doe",
				Email:    "john@doe.com",
				Password: "Pa$$w0rd",
			}),
		).
		Response(
			apiduck.JSONResponse(http.StatusCreated, struct {
				User users.User `json:"user"`
			}{}).Example(map[string]any{
				"user": users.User{
					Id:        12,
					Name:      "John Doe",
					Email:     "john@doe.com",
					Role:      users.RoleEmployee,
					IsActive:  true,
					CreatedAt: time.Now().Format(time.DateOnly),
				},
			}),
		).
		Response(
			apiduck.JSONResponse(http.StatusBadRequest, errorEnvelope{}).Example(errorEnvelope{
				Code:  ErrorCodeBadRequest,
				Error: "invalid password",
			}),
		).
		Response(
			apiduck.JSONResponse(http.StatusConflict, errorEnvelope{}).Example(errorEnvelope{
				Code:  ErrorCodeConflict,
				Error: "email is already in use",
			}),
		).
		Response(
			apiduck.JSONResponse(http.StatusInternalServerError, errorEnvelope{}).Example(errorEnvelope{
				Code:  ErrorCodeInternal,
				Error: "something went wrong",
			}),
		)

	authResource.Post("/v1/auth/login", "Log ind", "Log ind med email og password").
		Body(
			apiduck.JSONBody(users.LoginUserRequest{}).Example(users.LoginUserRequest{
				Email:    "john@doe.com",
				Password: "Pa$$w0rd",
			}),
		).
		Response(
			apiduck.JSONResponse(http.StatusCreated, struct {
				Session sessions.Session `json:"session"`
			}{}).Example(map[string]any{
				"session": sessions.Session{
					Token:     "0M86PHQQDG72M1OGLOIMULIDQ9ILN6V3",
					UserId:    12,
					ExpiresAt: time.Now().Add(7 * 24 * time.Hour).Format(time.DateTime),
					CreatedAt: time.Now().Format(time.DateTime),
					UpdatedAt: time.Now().Format(time.DateTime),
				},
			}),
		).
		Response(
			apiduck.JSONResponse(http.StatusUnauthorized, errorEnvelope{}).Example(errorEnvelope{
				Code:  ErrorCodeUnauthorized,
				Error: "invalid credentials",
			}),
		).
		Response(
			apiduck.JSONResponse(http.StatusInternalServerError, errorEnvelope{}).Example(errorEnvelope{
				Code:  ErrorCodeInternal,
				Error: "something went wrong",
			}),
		)

	meResource := docs.AddResource("Me", "Tidsregistreringer og kategorier")

	meResource.Get("/v1/me/categories", "Hent followed kategorier", "Henter kategorier som brugeren har valgt at follow").
		Security("(bearer-token-for-users)").
		Response(
			apiduck.JSONResponse(http.StatusOK, struct {
				Categories []categories.Category `json:"categories"`
			}{}).Example(map[string]any{
				"categories": []categories.Category{
					{
						Id:        2,
						Title:     "Support",
						RootTitle: "Support",
					},
					{
						Id:        1,
						Title:     "Software Udvikling",
						RootTitle: "Software Udvikling",
					},
					{
						Id:        3,
						Title:     "Standardløsning",
						RootTitle: "Software Udvikling",
					},
				},
			}),
		).
		Response(
			apiduck.JSONResponse(http.StatusInternalServerError, errorEnvelope{}).Example(errorEnvelope{
				Code:  ErrorCodeInternal,
				Error: "something went wrong",
			}),
		)

	meResource.Post("/v1/me/categories", "Opret ny kategori", "Opret en ny kateogri som enten root-kategori eller som en child-kategori").
		Security("(bearer-token-for-users)").
		Body(apiduck.JSONBody(categories.CreateCategoryInput{}).Example(categories.CreateCategoryInput{
			Title:    "Lagerløsning",
			ParentId: nil,
		})).
		Response(apiduck.JSONResponse(http.StatusCreated, struct {
			Category categories.Category `json:"category"`
		}{}).Example(struct {
			Category categories.Category `json:"category"`
		}{
			Category: categories.Category{
				Id:        42,
				Title:     "Lagerløsning",
				RootTitle: "",
			},
		})).
		Response(
			apiduck.JSONResponse(http.StatusBadRequest, errorEnvelope{}).Example(errorEnvelope{
				Code:  ErrorCodeBadRequest,
				Error: "invalid body",
			}),
		).
		Response(
			apiduck.JSONResponse(http.StatusInternalServerError, errorEnvelope{}).Example(errorEnvelope{
				Code:  ErrorCodeInternal,
				Error: "something went wrong",
			}),
		)

	meResource.Put("/v1/me/categories/{id}", "Opdater en kategori", "Opdater en kategoris titel").
		Security("(bearer-token-for-users)").
		PathParams(apiduck.PathParam("id", "Kategori id").Example(42)).
		Body(apiduck.JSONBody(categories.UpdateCategoryInput{}).Example(categories.UpdateCategoryInput{
			Title: "Ny Lagerløsning for kudne",
		})).
		Response(apiduck.JSONResponse(http.StatusCreated, struct {
			Category categories.Category `json:"category"`
		}{}).Example(struct {
			Category categories.Category `json:"category"`
		}{
			Category: categories.Category{
				Id:        42,
				Title:     "Ny Lagerløsning for kunde",
				RootTitle: "",
			},
		})).
		Response(
			apiduck.JSONResponse(http.StatusBadRequest, errorEnvelope{}).Example(errorEnvelope{
				Code:  ErrorCodeBadRequest,
				Error: "invalid body",
			}),
		).
		Response(
			apiduck.JSONResponse(http.StatusInternalServerError, errorEnvelope{}).Example(errorEnvelope{
				Code:  ErrorCodeInternal,
				Error: "something went wrong",
			}),
		)

	meResource.Put("/v1/me/categories/{id}/toggle", "Spær eller åben en kategori", "Ændrer status på en kategori mellem spærret og åbnet").
		Security("(bearer-token-for-users)").
		PathParams(apiduck.PathParam("id", "Kategori id").Example(42)).
		Response(apiduck.JSONResponse(http.StatusNoContent, nil).Description("Kategori spærret/åbnet")).
		Response(
			apiduck.JSONResponse(http.StatusBadRequest, errorEnvelope{}).Example(errorEnvelope{
				Code:  ErrorCodeBadRequest,
				Error: "invalid body",
			}),
		).
		Response(
			apiduck.JSONResponse(http.StatusInternalServerError, errorEnvelope{}).Example(errorEnvelope{
				Code:  ErrorCodeInternal,
				Error: "something went wrong",
			}),
		)

	meResource.Get("/v1/me/categories/all", "Hent alle kategorier", "Henter alle kategorier grupperet under deres \"parent\"-kategori").
		Security("(bearer-token-for-users)").
		Response(
			apiduck.JSONResponse(http.StatusOK, struct {
				Categories []categories.CategoryTree `json:"categories"`
			}{}).Example(map[string]any{
				"categories": []categories.CategoryTree{
					{
						Id:         2,
						ParentId:   nil,
						Title:      "Software Udvikling",
						IsRetired:  false,
						IsFollowed: true,
						Children: []*categories.CategoryTree{
							{
								Id:         4,
								ParentId:   ptr(int64(2)),
								Title:      "Skræddersyet",
								IsRetired:  false,
								IsFollowed: false,
								Children:   []*categories.CategoryTree{},
							},
						},
					},
				},
			}),
		).
		Response(
			apiduck.JSONResponse(http.StatusInternalServerError, errorEnvelope{}).Example(errorEnvelope{
				Code:  ErrorCodeInternal,
				Error: "something went wrong",
			}),
		)

	meResource.Put("/v1/me/categories/{id}/follow", "Follow en kategori", "Follow en kategori så brugeren kan lave tids registreringer på den kategori").
		Security("(bearer-token-for-users)").
		PathParams(
			apiduck.PathParam("id", "Kategori id").Example(42),
		).
		Response(
			apiduck.JSONResponse(http.StatusNoContent, nil).Description("Kategori er followed"),
		).
		Response(
			apiduck.JSONResponse(http.StatusConflict, errorEnvelope{}).Example(errorEnvelope{
				Code:  ErrorCodeConflict,
				Error: "already following category",
			}),
		).
		Response(
			apiduck.JSONResponse(http.StatusInternalServerError, errorEnvelope{}).Example(errorEnvelope{
				Code:  ErrorCodeInternal,
				Error: "something went wrong",
			}),
		)

	meResource.Put("/v1/me/categories/{id}/unfollow", "Unfollow en kategori", "Unfollow en kategori så brugeren ikke kan lave tids registreringer på den kategori længere").
		Security("(bearer-token-for-users)").
		PathParams(
			apiduck.PathParam("id", "Kategori id").Example(42),
		).
		Response(
			apiduck.JSONResponse(http.StatusNoContent, nil).Description("Kategori er unfollowed"),
		).
		Response(
			apiduck.JSONResponse(http.StatusBadRequest, errorEnvelope{}).Example(errorEnvelope{
				Code:  ErrorCodeBadRequest,
				Error: "invalied catogory id",
			}),
		).
		Response(
			apiduck.JSONResponse(http.StatusNotFound, errorEnvelope{}).Example(errorEnvelope{
				Code:  ErrorCodeNotFound,
				Error: "category is not followed",
			}),
		).
		Response(
			apiduck.JSONResponse(http.StatusInternalServerError, errorEnvelope{}).Example(errorEnvelope{
				Code:  ErrorCodeInternal,
				Error: "something went wrong",
			}),
		)

	meResource.Post("/v1/me/time_entries", "Opret ny tidsregistrering", "Opret en ny tidsregistrering for en given dato").
		Security("(bearer-token-for-users)").
		Body(
			apiduck.JSONBody(time_entries.RegisterTimeEntryInput{}).Example(time_entries.RegisterTimeEntryInput{
				CategoryId:  42,
				Date:        time.Now().Format(time.DateOnly),
				Duration:    types.Duration{Duration: 2*time.Hour + 30*time.Minute},
				Description: "Ny feature implementeret",
			}),
		).
		Response(
			apiduck.JSONResponse(http.StatusCreated, struct {
				TimeEntry time_entries.TimeEntry `json:"timeEntry"`
			}{}).Example(map[string]any{
				"timeEntry": time_entries.TimeEntry{
					Id:          2,
					CategoryId:  42,
					Category:    "Ny app idé",
					UserId:      32,
					Date:        time.Now().Format(time.DateOnly),
					Duration:    types.Duration{Duration: 2*time.Hour + 30*time.Minute},
					Description: "Ny feature implementeret",
				},
			}),
		).
		Response(
			apiduck.JSONResponse(http.StatusBadRequest, errorEnvelope{}).Example(errorEnvelope{
				Code:  ErrorCodeBadRequest,
				Error: "invalid body",
			}),
		).
		Response(
			apiduck.JSONResponse(http.StatusInternalServerError, errorEnvelope{}).Example(errorEnvelope{
				Code:  ErrorCodeInternal,
				Error: "something went wrong",
			}),
		)

	meResource.Put("/v1/me/time_entries/{id}", "Opdater en tidsregistrering", "Opdater en tidsregistrering").
		Security("(bearer-token-for-users)").
		PathParams(
			apiduck.PathParam("id", "Tidsregistrerings id").Example(2),
		).
		Body(
			apiduck.JSONBody(time_entries.UpdateTimeEntryInput{}).Example(time_entries.UpdateTimeEntryInput{
				Duration:    types.Duration{Duration: 3*time.Hour + 40*time.Minute},
				Description: "Ny feature implementeret + unit tests",
			}),
		).
		Response(
			apiduck.JSONResponse(http.StatusCreated, struct {
				TimeEntry time_entries.TimeEntry `json:"timeEntry"`
			}{}).Example(map[string]any{
				"timeEntry": time_entries.TimeEntry{
					Id:          2,
					CategoryId:  42,
					Category:    "Ny app idé",
					UserId:      32,
					Date:        time.Now().Format(time.DateOnly),
					Duration:    types.Duration{Duration: 3*time.Hour + 40*time.Minute},
					Description: "Ny feature implementeret + unit tests",
				},
			}),
		).
		Response(
			apiduck.JSONResponse(http.StatusBadRequest, errorEnvelope{}).Example(errorEnvelope{
				Code:  ErrorCodeBadRequest,
				Error: "invalid body",
			}),
		).
		Response(
			apiduck.JSONResponse(http.StatusInternalServerError, errorEnvelope{}).Example(errorEnvelope{
				Code:  ErrorCodeInternal,
				Error: "something went wrong",
			}),
		)

	meResource.Delete("/v1/me/time_entries/{id}", "Slet en tidsregistrering", "Slet en tidsregistrering").
		Security("(bearer-token-for-users)").
		PathParams(
			apiduck.PathParam("id", "Tidsregistrerings id").Example(10),
		).
		Response(
			apiduck.JSONResponse(http.StatusNoContent, nil).Description("Tidsregistrering blev slettet"),
		).
		Response(
			apiduck.JSONResponse(http.StatusBadRequest, errorEnvelope{}).Example(errorEnvelope{
				Code:  ErrorCodeBadRequest,
				Error: "invalied time entry id",
			}),
		).
		Response(
			apiduck.JSONResponse(http.StatusInternalServerError, errorEnvelope{}).Example(errorEnvelope{
				Code:  ErrorCodeInternal,
				Error: "something went wrong",
			}),
		)

	meResource.Get("/v1/me/time_entries/day/{date}", "Hent tidsregistreringer for dato", "Hent tidsregistreringer for dato med samlet antal tid brugt").
		Security("(bearer-token-for-users)").
		PathParams(
			apiduck.PathParam("date", "Dato").Example(time.Now().Format(time.DateOnly)),
		).
		Response(
			apiduck.JSONResponse(
				http.StatusOK,
				struct {
					Summary time_entries.SummaryDay `json:"summary"`
				}{}).
				Example(map[string]any{
					"summary": time_entries.SummaryDay{
						Date:       time.Now().Format(time.DateOnly),
						Weekday:    strings.ToLower(time.Now().Weekday().String()),
						TotalHours: types.Duration{Duration: 4 * time.Hour},
						MaxHours:   types.Duration{Duration: 7 * time.Hour},
						TimeEntries: []time_entries.TimeEntry{
							{
								Id:          1,
								CategoryId:  3,
								Category:    "Support",
								UserId:      23,
								Date:        time.Now().Format(time.DateOnly),
								Duration:    types.Duration{Duration: 4 * time.Hour},
								Description: "Oprettet nye brugere for kunde",
							},
						},
					},
				}),
		).
		Response(
			apiduck.JSONResponse(http.StatusBadRequest, errorEnvelope{}).Example(errorEnvelope{
				Code:  ErrorCodeBadRequest,
				Error: "invalid body",
			}),
		).
		Response(
			apiduck.JSONResponse(http.StatusInternalServerError, errorEnvelope{}).Example(errorEnvelope{
				Code:  ErrorCodeInternal,
				Error: "something went wrong",
			}),
		)

	meResource.Get("/v1/me/time_entries/month/{year-month}", "Hent tidsregistreringer for måned", "Hent tidsregistreringer for måned med samlet antal tid brugt").
		Security("(bearer-token-for-users)").
		PathParams(
			apiduck.PathParam("year-month", "År og måned").Example(time.Now().Format("2006-01")),
		).
		Response(
			apiduck.JSONResponse(
				http.StatusOK,
				struct {
					Summary time_entries.SummaryMonth `json:"summary"`
				}{}).
				Example(map[string]any{
					"summary": time_entries.SummaryMonth{
						Month:      strings.ToLower(time.Now().Month().String()),
						TotalHours: types.Duration{Duration: 24 * 14 * time.Hour},
						MaxHours:   types.Duration{Duration: 24 * 22 * time.Hour},
						Days: []time_entries.SummaryDay{
							{
								Date:       time.Now().Format(time.DateOnly),
								Weekday:    strings.ToLower(time.Now().Weekday().String()),
								TotalHours: types.Duration{Duration: 4 * time.Hour},
								MaxHours:   types.Duration{Duration: 7 * time.Hour},
								TimeEntries: []time_entries.TimeEntry{
									{
										Id:          1,
										CategoryId:  3,
										Category:    "Support",
										UserId:      23,
										Date:        time.Now().Format(time.DateOnly),
										Duration:    types.Duration{Duration: 4 * time.Hour},
										Description: "Oprettet nye brugere for kunde",
									},
								},
							},
						},
					},
				}),
		).
		Response(
			apiduck.JSONResponse(http.StatusBadRequest, errorEnvelope{}).Example(errorEnvelope{
				Code:  ErrorCodeBadRequest,
				Error: "invalid body",
			}),
		).
		Response(
			apiduck.JSONResponse(http.StatusInternalServerError, errorEnvelope{}).Example(errorEnvelope{
				Code:  ErrorCodeInternal,
				Error: "something went wrong",
			}),
		)

	return docs
}
