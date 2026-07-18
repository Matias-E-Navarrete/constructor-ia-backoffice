package httpserver

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"rimu/backend/internal/platform/config"
	"rimu/backend/internal/platform/db"
	"rimu/backend/internal/platform/httpmiddleware"

	ffapplication "rimu/backend/internal/platform/featureflags/application"
	ffinfrastructure "rimu/backend/internal/platform/featureflags/infrastructure"

	userapplication "rimu/backend/internal/user/application"
	userinfrastructure "rimu/backend/internal/user/infrastructure"
	userhttp "rimu/backend/internal/user/interfaces/http"

	groupsapplication "rimu/backend/internal/groups/application"
	groupsinfrastructure "rimu/backend/internal/groups/infrastructure"
	groupshttp "rimu/backend/internal/groups/interfaces/http"

	habitsapplication "rimu/backend/internal/habits/application"
	habitsinfrastructure "rimu/backend/internal/habits/infrastructure"
	habitshttp "rimu/backend/internal/habits/interfaces/http"

	workoutsapplication "rimu/backend/internal/workouts/application"
	workoutsinfrastructure "rimu/backend/internal/workouts/infrastructure"
	workoutshttp "rimu/backend/internal/workouts/interfaces/http"

	financeapplication "rimu/backend/internal/finance/application"
	financeinfrastructure "rimu/backend/internal/finance/infrastructure"
	financehttp "rimu/backend/internal/finance/interfaces/http"

	notesapplication "rimu/backend/internal/notes/application"
	notesinfrastructure "rimu/backend/internal/notes/infrastructure"
	noteshttp "rimu/backend/internal/notes/interfaces/http"

	adminapplication "rimu/backend/internal/admin/application"
	admininfrastructure "rimu/backend/internal/admin/infrastructure"
	adminhttp "rimu/backend/internal/admin/interfaces/http"

	tasksapplication "rimu/backend/internal/tasks/application"
	tasksinfrastructure "rimu/backend/internal/tasks/infrastructure"
	taskshttp "rimu/backend/internal/tasks/interfaces/http"

	inboxapplication "rimu/backend/internal/inbox/application"
	inboxinfrastructure "rimu/backend/internal/inbox/infrastructure"
	inboxhttp "rimu/backend/internal/inbox/interfaces/http"
)

// Build is the composition root logic shared by cmd/api/main.go and the
// backend's own e2e test suite: connect, migrate, wire every bounded
// context, and return the fully-mounted router.
func Build(ctx context.Context, cfg config.Config) (*chi.Mux, *pgxpool.Pool, error) {
	pool, err := db.Connect(ctx, cfg.DatabaseURL)
	if err != nil {
		return nil, nil, err
	}

	if err := db.Migrate(ctx, pool); err != nil {
		pool.Close()
		return nil, nil, err
	}

	flagsRepo := ffinfrastructure.NewPostgresRepository(pool)
	flags := ffapplication.NewRegistry(flagsRepo)
	if err := flags.Load(ctx); err != nil {
		pool.Close()
		return nil, nil, err
	}

	userRepo := userinfrastructure.NewPostgresRepository(pool)
	groupsRepo := groupsinfrastructure.NewPostgresRepository(pool)
	habitsRepo := habitsinfrastructure.NewPostgresRepository(pool)
	workoutsRepo := workoutsinfrastructure.NewPostgresRepository(pool)
	financeRepo := financeinfrastructure.NewPostgresRepository(pool)
	notesRepo := notesinfrastructure.NewPostgresRepository(pool)
	roadmapRepo := admininfrastructure.NewPostgresRoadmapRepository(pool)
	tasksRepo := tasksinfrastructure.NewPostgresRepository(pool)
	inboxRepo := inboxinfrastructure.NewPostgresRepository(pool)
	accountsRepo := financeinfrastructure.NewPostgresAccountRepository(pool)

	userHandler := &userhttp.Handler{
		Register:           &userapplication.RegisterUser{Repo: userRepo, IsAdminEmail: cfg.IsAdminEmail},
		Login:              &userapplication.AuthenticateUser{Repo: userRepo, JWTSecret: cfg.JWTSecret},
		Upgrade:            &userapplication.UpgradePlan{Repo: userRepo},
		UpdatePlannerHours: &userapplication.UpdatePlannerHours{Repo: userRepo},
	}

	groupsHandler := &groupshttp.Handler{
		CreateFamily:   &groupsapplication.CreateFamilyGroup{Repo: groupsRepo},
		CreateCoaching: &groupsapplication.CreateCoachingGroup{Repo: groupsRepo},
		InviteMember:   &groupsapplication.InviteMember{Groups: groupsRepo, Users: userRepo},
		RemoveMember:   &groupsapplication.RemoveMember{Repo: groupsRepo},
		ListMyGroups:   &groupsapplication.ListMyGroups{Repo: groupsRepo},
		GetGroup:       &groupsapplication.GetGroup{Repo: groupsRepo},
	}

	habitsHandler := &habitshttp.Handler{
		Create:   &habitsapplication.CreateHabit{Repo: habitsRepo},
		Update:   &habitsapplication.UpdateHabit{Repo: habitsRepo},
		Delete:   &habitsapplication.DeleteHabit{Repo: habitsRepo},
		CheckIn:  &habitsapplication.CheckInHabit{Repo: habitsRepo},
		List:     &habitsapplication.ListHabits{Repo: habitsRepo, Groups: groupsRepo},
		GetStats: &habitsapplication.GetHabitStats{Repo: habitsRepo},
	}

	workoutsHandler := &workoutshttp.Handler{
		LogSession:   &workoutsapplication.LogSession{Repo: workoutsRepo},
		ListSessions: &workoutsapplication.ListSessions{Repo: workoutsRepo, Groups: groupsRepo},
		GetSession:   &workoutsapplication.GetSession{Repo: workoutsRepo},
		Delete:       &workoutsapplication.DeleteSession{Repo: workoutsRepo},
		GetProgress:  &workoutsapplication.GetExerciseProgress{Repo: workoutsRepo},
	}

	financeHandler := &financehttp.Handler{
		Record:        &financeapplication.RecordTransaction{Repo: financeRepo, Groups: groupsRepo},
		List:          &financeapplication.ListTransactions{Repo: financeRepo, Groups: groupsRepo},
		Delete:        &financeapplication.DeleteTransaction{Repo: financeRepo},
		GetSummary:    &financeapplication.GetFinanceSummary{Repo: financeRepo},
		ExportCSV:     &financeapplication.ExportTransactions{Repo: financeRepo},
		CreateAccount: &financeapplication.CreateAccount{Repo: accountsRepo},
		ListAccounts:  &financeapplication.ListAccounts{Repo: accountsRepo},
		Convert:       &financeapplication.ConvertCurrency{},
		GetUpcoming:   &financeapplication.GetUpcomingBills{Repo: financeRepo},
		ExportPDF:     &financeapplication.ExportPDF{Repo: financeRepo},
	}

	notesHandler := &noteshttp.Handler{
		Create: &notesapplication.CreateNote{Repo: notesRepo},
		Update: &notesapplication.UpdateNote{Repo: notesRepo},
		Get:    &notesapplication.GetNote{Repo: notesRepo},
		List:   &notesapplication.ListNotes{Repo: notesRepo},
		Delete: &notesapplication.DeleteNote{Repo: notesRepo},
		GetGraph: &notesapplication.GetVaultGraph{
			Notes: notesRepo, Habits: habitsRepo, Workouts: workoutsRepo, Finance: financeRepo,
		},
		ExportVault: &notesapplication.ExportVault{
			Notes: notesRepo, Habits: habitsRepo, Workouts: workoutsRepo, Finance: financeRepo,
		},
	}

	tasksHandler := &taskshttp.Handler{
		Create:      &tasksapplication.CreateTask{Repo: tasksRepo},
		Update:      &tasksapplication.UpdateTask{Repo: tasksRepo},
		Delete:      &tasksapplication.DeleteTask{Repo: tasksRepo},
		List:        &tasksapplication.ListTasks{Repo: tasksRepo},
		Complete:    &tasksapplication.CompleteTask{Repo: tasksRepo},
		Uncomplete:  &tasksapplication.UncompleteTask{Repo: tasksRepo},
		SetStatus:   &tasksapplication.SetStatus{Repo: tasksRepo},
		SetQuadrant: &tasksapplication.SetQuadrant{Repo: tasksRepo},
		Reorder:     &tasksapplication.ReorderTasks{Repo: tasksRepo},
		Schedule:    &tasksapplication.ScheduleTask{Repo: tasksRepo},
	}

	inboxHandler := &inboxhttp.Handler{
		Capture:   &inboxapplication.CaptureItem{Repo: inboxRepo},
		List:      &inboxapplication.ListInbox{Repo: inboxRepo},
		SetPinned: &inboxapplication.SetPinned{Repo: inboxRepo},
		Delete:    &inboxapplication.DeleteItem{Repo: inboxRepo},
	}

	adminHandler := &adminhttp.Handler{
		ListRoadmap:         &adminapplication.ListRoadmap{Repo: roadmapRepo},
		CreateRoadmapItem:   &adminapplication.CreateRoadmapItem{Repo: roadmapRepo},
		UpdateRoadmapStatus: &adminapplication.UpdateRoadmapStatus{Repo: roadmapRepo},
		ListFeatureFlags:    &adminapplication.ListFeatureFlags{Registry: flags},
		ToggleFeatureFlag:   &adminapplication.ToggleFeatureFlag{Registry: flags},
		ListUsers:           &adminapplication.ListUsers{Repo: userRepo},
		SetUserPlan:         &adminapplication.SetUserPlan{Repo: userRepo},
		SetUserRole:         &adminapplication.SetUserRole{Repo: userRepo},
		RunBackendTests:     &adminapplication.RunBackendTests{WorkDir: cfg.BackendDir},
		RunFrontendTests:    &adminapplication.RunFrontendTests{WorkDir: cfg.FrontendDir},
	}

	requireAuth := httpmiddleware.RequireAuth(cfg.JWTSecret, userRepo)
	requirePro := httpmiddleware.RequirePro
	requireAdmin := httpmiddleware.RequireAdmin
	feature := func(key string) func(http.Handler) http.Handler {
		return httpmiddleware.RequireFeatureEnabled(flags, key)
	}

	router := New()
	router.Route("/api", func(api chi.Router) {
		userhttp.Mount(api, userHandler, requireAuth)
		groupshttp.Mount(api, groupsHandler, requireAuth, feature("groups"))

		habitshttp.Mount(api, habitsHandler, habitshttp.Middlewares{
			RequireAuth:       requireAuth,
			RequireModuleFlag: feature("habits"),
			RequireStatsFlag:  feature("habits.stats"),
			RequirePro:        requirePro,
		})

		workoutshttp.Mount(api, workoutsHandler, workoutshttp.Middlewares{
			RequireAuth:         requireAuth,
			RequireModuleFlag:   feature("workouts"),
			RequireProgressFlag: feature("workouts.progress"),
			RequirePro:          requirePro,
		})

		financehttp.Mount(api, financeHandler, financehttp.Middlewares{
			RequireAuth:        requireAuth,
			RequireModuleFlag:  feature("finance"),
			RequireSummaryFlag: feature("finance.summary"),
			RequireExportFlag:  feature("finance.export"),
			RequirePro:         requirePro,
		})

		noteshttp.Mount(api, notesHandler, noteshttp.Middlewares{
			RequireAuth:       requireAuth,
			RequireModuleFlag: feature("notes"),
			RequireGraphFlag:  feature("notes.graph"),
			RequireExportFlag: feature("notes.export"),
			RequirePro:        requirePro,
		})

		taskshttp.Mount(api, tasksHandler, requireAuth, feature("tasks"))
		inboxhttp.Mount(api, inboxHandler, requireAuth, feature("inbox"))

		adminhttp.Mount(api, adminHandler, requireAuth, requireAdmin)
	})

	return router, pool, nil
}
