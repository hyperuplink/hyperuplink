package activity

import (
	"xn--gckvb8fzb.com/hyperuplink/models/activity"
	"xn--gckvb8fzb.com/hyperuplink/models/setting"
	"xn--gckvb8fzb.com/hyperuplink/runtime"
	repoSetting "xn--gckvb8fzb.com/hyperuplink/services/repositories/setting"
)

const CleanupAdminLogID string = "admin_log_cleanup"

const CleanupAdminLogSpec string = "0 0 * * *"

func AdminLogRetentionDays(rt *runtime.Runtime) int {
	settingSystem, err := repoSetting.GetByID[setting.System](
		rt.Repositories.Setting,
		"system",
	)
	if err != nil {
		rt.Error("error", err)
		return setting.DEFAULT_ADMIN_LOG_RETENTION_DAYS
	}

	return settingSystem.JSONValue.GetAdminLogRetentionDays()
}

func CleanupAdminLog(rt *runtime.Runtime) (err error) {
	days := AdminLogRetentionDays(rt)

	deleted, err := rt.Repositories.Activity.DeleteOlderThanDays(
		activity.AdminKinds(),
		days,
	)
	if err != nil {
		return err
	}

	rt.Info("cleanup", "admin_log", "days", days, "deleted", deleted)

	return nil
}
