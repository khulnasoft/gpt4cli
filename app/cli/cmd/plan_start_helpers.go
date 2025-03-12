package cmd

import (
	"fmt"
	"gpt4cli-cli/api"
	"gpt4cli-cli/lib"
	"gpt4cli-cli/term"
	shared "gpt4cli-shared"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	// Tier flags
	noAuto    bool
	basicAuto bool
	plusAuto  bool
	semiAuto  bool
	fullAuto  bool

	// Type flags
	dailyModels  bool
	strongModels bool
	ossModels    bool
	cheapModels  bool
)

func AddNewPlanFlags(cmd *cobra.Command) {
	// Add tier flags
	cmd.Flags().BoolVar(&noAuto, "no-auto", false, shared.AutoModeDescriptions[shared.AutoModeNone])
	cmd.Flags().BoolVar(&basicAuto, "basic", false, shared.AutoModeDescriptions[shared.AutoModeBasic])
	cmd.Flags().BoolVar(&plusAuto, "plus", false, shared.AutoModeDescriptions[shared.AutoModePlus])
	cmd.Flags().BoolVar(&semiAuto, "semi", false, shared.AutoModeDescriptions[shared.AutoModeSemi])
	cmd.Flags().BoolVar(&fullAuto, "full", false, shared.AutoModeDescriptions[shared.AutoModeFull])

	// Add type flags
	cmd.Flags().BoolVar(&dailyModels, "daily", false, shared.DailyDriverModelPack.Description)
	cmd.Flags().BoolVar(&strongModels, "strong", false, shared.StrongModelPack.Description)
	cmd.Flags().BoolVar(&cheapModels, "cheap", false, shared.CheapModelPack.Description)
	cmd.Flags().BoolVar(&ossModels, "oss", false, shared.OSSModelPack.Description)
}

func resolveAutoMode(config *shared.PlanConfig) (bool, *shared.PlanConfig) {
	didUpdate, updatedConfig, _ := resolveAutoModeWithArgs(config, false)
	return didUpdate, updatedConfig
}

func resolveAutoModeSilent(config *shared.PlanConfig) (bool, *shared.PlanConfig, func()) {
	return resolveAutoModeWithArgs(config, true)
}

func resolveAutoModeWithArgs(config *shared.PlanConfig, silent bool) (bool, *shared.PlanConfig, func()) {
	currentAutoMode := config.AutoMode
	var toSetAutoMode shared.AutoModeType
	if noAuto {
		toSetAutoMode = shared.AutoModeNone
	} else if basicAuto {
		toSetAutoMode = shared.AutoModeBasic
	} else if plusAuto {
		toSetAutoMode = shared.AutoModePlus
	} else if semiAuto {
		toSetAutoMode = shared.AutoModeSemi
	} else if fullAuto {
		toSetAutoMode = shared.AutoModeFull
	}

	if toSetAutoMode != "" && toSetAutoMode != currentAutoMode {
		if !silent {
			term.StartSpinner("")
		}
		_, updatedConfig := updateConfig([]string{"auto-mode", string(toSetAutoMode)}, config)
		apiErr := api.Client.UpdatePlanConfig(lib.CurrentPlanId, shared.UpdatePlanConfigRequest{
			Config: updatedConfig,
		})
		if !silent {
			term.StopSpinner()
		}

		if apiErr != nil {
			term.OutputErrorAndExit("Error updating config auto-mode: %v", apiErr)
		}

		fn := func() {
			fmt.Printf("🚀 Set %s to %s\n", color.New(color.Bold).Sprint("auto-mode"), color.New(color.Bold, term.ColorHiMagenta).Sprint(toSetAutoMode))
		}

		if !silent {
			fn()
			return true, updatedConfig, fn
		}

		return true, updatedConfig, fn
	}

	return false, config, nil
}

func resolveModelPack() {
	resolveModelPackWithArgs(nil, false)
}

func resolveModelPackSilent(settings *shared.PlanSettings) (*shared.PlanSettings, func()) {
	return resolveModelPackWithArgs(settings, true)
}

func resolveModelPackWithArgs(settings *shared.PlanSettings, silent bool) (*shared.PlanSettings, func()) {

	var originalSettings *shared.PlanSettings
	var apiErr *shared.ApiError
	if settings == nil {
		if !silent {
			term.StartSpinner("")
		}
		originalSettings, apiErr = api.Client.GetSettings(lib.CurrentPlanId, lib.CurrentBranch)
	} else {
		originalSettings = settings
	}

	if apiErr != nil {
		term.OutputErrorAndExit("Error getting current settings: %v", apiErr)
		return nil, nil
	}

	var packName string

	if ossModels {
		packName = shared.OSSModelPack.Name
	} else if strongModels {
		packName = shared.StrongModelPack.Name
	} else if cheapModels {
		packName = shared.CheapModelPack.Name
	} else {
		packName = shared.DailyDriverModelPack.Name
	}

	if packName != originalSettings.ModelPack.Name {
		if !silent {
			term.StartSpinner("")
		}
		updatedSettings := updateModelSettings([]string{packName}, originalSettings)
		_, apiErr = api.Client.UpdateSettings(lib.CurrentPlanId, lib.CurrentBranch, shared.UpdateSettingsRequest{
			Settings: updatedSettings,
		})
		if !silent {
			term.StopSpinner()
		}

		if apiErr != nil {
			term.OutputErrorAndExit("Error setting model pack: %v", apiErr)
			return nil, nil
		}

		fn := func() {
			fmt.Printf("🧠 Set model pack to %s\n", color.New(color.Bold, term.ColorHiMagenta).Sprint(packName))
		}

		if !silent {
			fn()
			return updatedSettings, fn
		}

		return updatedSettings, fn
	} else {
		if !silent {
			term.StopSpinner()
		}
		fn := func() {
			fmt.Printf("🧠 Using model pack %s\n", color.New(color.Bold, term.ColorHiMagenta).Sprint(packName))
		}

		if !silent {
			fn()
			return originalSettings, fn
		}

		return originalSettings, fn
	}
}
