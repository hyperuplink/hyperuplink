package activity

import (
	"context"
	"xn--gckvb8fzb.com/glides/runtime"
	gh "xn--gckvb8fzb.com/hyperuplink/helpers"
	"xn--gckvb8fzb.com/hyperuplink/models/activity"
	"xn--gckvb8fzb.com/hyperuplink/models/setting"
	repoSetting "xn--gckvb8fzb.com/hyperuplink/services/repositories/setting"
)

const CleanupAdminLogID string = "admin_log_cleanup"

const CleanupAdminLogSpec string = "0 0 * * *"

func AdminLogRetentionDays(ctx context.Context, rt *runtime.Runtime) int {
	settingSystem, err := repoSetting.GetByID[setting.System](ctx,
		gh.Repositories(rt).Setting,
		"system",
	)
	if err != nil {
		rt.Error("error", err)
		return setting.DEFAULT_ADMIN_LOG_RETENTION_DAYS
	}

	return settingSystem.JSONValue.GetAdminLogRetentionDays()
}

func CleanupAdminLog(ctx context.Context, rt *runtime.Runtime) (err error) {
	days := AdminLogRetentionDays(ctx, rt)

	deleted, err := gh.Repositories(rt).Activity.DeleteOlderThanDays(ctx,
		activity.AdminKinds(),
		days,
	)
	if err != nil {
		return err
	}

	rt.Info("cleanup", "admin_log", "days", days, "deleted", deleted)

	return nil
}
