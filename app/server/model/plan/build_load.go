package plan

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"gpt4cli-server/db"
	"gpt4cli-server/syntax"
	"gpt4cli-server/types"

	"github.com/khulnasoft/gpt4cli/shared"
)

func (state *activeBuildStreamState) loadPendingBuilds() (map[string][]*types.ActiveBuild, error) {
	clients := state.clients
	plan := state.plan
	branch := state.branch
	auth := state.auth

	active, err := activatePlan(clients, plan, branch, auth, "", true, false)

	if err != nil {
		log.Printf("Error activating plan: %v\n", err)
	}

	repoLockId, err := db.LockRepo(
		db.LockRepoParams{
			OrgId:    auth.OrgId,
			UserId:   auth.User.Id,
			PlanId:   plan.Id,
			Branch:   branch,
			Scope:    db.LockScopeRead,
			Ctx:      active.Ctx,
			CancelFn: active.CancelFn,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("error locking repo for build: %v", err)
	}

	var modelContext []*db.Context
	var pendingBuildsByPath map[string][]*types.ActiveBuild
	var settings *shared.PlanSettings

	err = func() error {
		defer func() {
			err := db.DeleteRepoLock(repoLockId)
			if err != nil {
				log.Printf("Error unlocking repo: %v\n", err)
			}
		}()

		errCh := make(chan error)

		go func() {
			res, err := db.GetPlanContexts(auth.OrgId, plan.Id, true, false)
			if err != nil {
				log.Printf("Error getting plan modelContext: %v\n", err)
				errCh <- fmt.Errorf("error getting plan modelContext: %v", err)
				return
			}
			modelContext = res

			errCh <- nil
		}()

		go func() {
			res, err := active.PendingBuildsByPath(auth.OrgId, auth.User.Id, nil)

			if err != nil {
				log.Printf("Error getting pending builds by path: %v\n", err)
				errCh <- fmt.Errorf("error getting pending builds by path: %v", err)
				return
			}

			pendingBuildsByPath = res

			errCh <- nil
		}()

		go func() {
			res, err := db.GetPlanSettings(plan, true)
			if err != nil {
				log.Printf("Error getting plan settings: %v\n", err)
				errCh <- fmt.Errorf("error getting plan settings: %v", err)
				return
			}

			settings = res
			errCh <- nil
		}()

		for i := 0; i < 3; i++ {
			err = <-errCh
			if err != nil {
				log.Printf("Error getting plan data: %v\n", err)
				return err
			}
		}
		return nil
	}()

	if err != nil {
		return nil, err
	}

	UpdateActivePlan(plan.Id, branch, func(ap *types.ActivePlan) {
		ap.Contexts = modelContext
		for _, context := range modelContext {
			if context.FilePath != "" {
				ap.ContextsByPath[context.FilePath] = context
			}
		}
	})

	state.modelContext = modelContext
	state.settings = settings

	return pendingBuildsByPath, nil
}

func (state *activeBuildStreamFileState) loadBuildFile(activeBuild *types.ActiveBuild) error {

	currentOrgId := state.currentOrgId
	currentUserId := state.currentUserId
	planId := state.plan.Id
	branch := state.branch
	filePath := state.filePath

	activePlan := GetActivePlan(planId, branch)

	if activePlan == nil {
		return fmt.Errorf("active plan not found")
	}

	convoMessageId := activeBuild.ReplyId

	ext := filepath.Ext(filePath)
	parser, lang, backupParser, backupLang := syntax.GetParserForExt(ext)

	if parser != nil {
		state.parser = parser
		state.language = lang
	} else if backupParser != nil {
		state.parser = backupParser
		state.language = backupLang
	}

	if state.parser != nil {
		validationRes, err := syntax.Validate(activePlan.Ctx, filePath, state.preBuildState)
		if err != nil {
			log.Printf(" error validating original file syntax: %v\n", err)
			return fmt.Errorf("error validating original file syntax: %v", err)
		}
		if validationRes.HasParser && !validationRes.TimedOut && !validationRes.Valid {
			state.preBuildStateSyntaxInvalid = true
		}
	}

	build := &db.PlanBuild{
		OrgId:          currentOrgId,
		PlanId:         planId,
		ConvoMessageId: convoMessageId,
		FilePath:       filePath,
	}
	err := db.StorePlanBuild(build)

	if err != nil {
		log.Printf("Error storing plan build: %v\n", err)
		UpdateActivePlan(activePlan.Id, activePlan.Branch, func(ap *types.ActivePlan) {
			ap.IsBuildingByPath[filePath] = false
		})
		activePlan.StreamDoneCh <- &shared.ApiError{
			Type:   shared.ApiErrorTypeOther,
			Status: http.StatusInternalServerError,
			Msg:    "Error storing plan build: " + err.Error(),
		}
		return err
	}

	var currentPlan *shared.CurrentPlanState
	var convo []*db.ConvoMessage

	log.Println("Locking repo for load build file")

	repoLockId, err := db.LockRepo(
		db.LockRepoParams{
			OrgId:       currentOrgId,
			UserId:      currentUserId,
			PlanId:      planId,
			Branch:      branch,
			PlanBuildId: build.Id,
			Scope:       db.LockScopeRead,
			Ctx:         activePlan.Ctx,
			CancelFn:    activePlan.CancelFn,
		},
	)
	if err != nil {
		log.Printf("Error locking repo for build file: %v\n", err)
		UpdateActivePlan(activePlan.Id, activePlan.Branch, func(ap *types.ActivePlan) {
			ap.IsBuildingByPath[filePath] = false
		})
		activePlan.StreamDoneCh <- &shared.ApiError{
			Type:   shared.ApiErrorTypeOther,
			Status: http.StatusInternalServerError,
			Msg:    "Error locking repo for build file: " + err.Error(),
		}
		return err
	}

	log.Println("Locked repo for load build file")

	err = func() error {
		defer func() {
			log.Printf("Unlocking repo for load build file")

			err := db.DeleteRepoLock(repoLockId)
			if err != nil {
				log.Printf("Error unlocking repo: %v\n", err)
			}
		}()

		errCh := make(chan error)

		go func() {
			log.Println("loadBuildFile - Getting current plan state")
			res, err := db.GetCurrentPlanState(db.CurrentPlanStateParams{
				OrgId:  currentOrgId,
				PlanId: planId,
			})
			if err != nil {
				log.Printf("Error getting current plan state: %v\n", err)
				UpdateActivePlan(activePlan.Id, activePlan.Branch, func(ap *types.ActivePlan) {
					ap.IsBuildingByPath[filePath] = false
				})
				activePlan.StreamDoneCh <- &shared.ApiError{
					Type:   shared.ApiErrorTypeOther,
					Status: http.StatusInternalServerError,
					Msg:    "Error getting current plan state: " + err.Error(),
				}
				errCh <- fmt.Errorf("error getting current plan state: %v", err)
				return
			}
			currentPlan = res

			log.Println("Got current plan state")
			errCh <- nil
		}()

		go func() {
			res, err := db.GetPlanConvo(currentOrgId, planId)
			if err != nil {
				log.Printf("Error getting plan convo: %v\n", err)
				errCh <- fmt.Errorf("error getting plan convo: %v", err)
				return
			}
			convo = res

			errCh <- nil
		}()

		for i := 0; i < 2; i++ {
			err = <-errCh
			if err != nil {
				log.Printf("Error getting plan data: %v\n", err)
				return err
			}
		}

		return nil

	}()

	if err != nil {
		return err
	}

	state.filePath = filePath
	state.convoMessageId = convoMessageId
	state.build = build
	state.currentPlanState = currentPlan
	state.convo = convo

	return nil

}
