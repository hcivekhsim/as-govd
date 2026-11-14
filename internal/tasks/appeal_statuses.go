package tasks

var statusGroups = map[string]StatusGroup{
	"initial": {
		Code:      "initial",
		Title:     "Начальные",
		BadgeType: "info",
	},
	"working": {
		Code:      "working",
		Title:     "Рабочие",
		BadgeType: "progress",
	},
	"intermediate": {
		Code:      "intermediate",
		Title:     "Промежуточные",
		BadgeType: "warning",
	},
	"final_positive": {
		Code:      "final_positive",
		Title:     "Финальные",
		BadgeType: "success",
	},
	"final_negative": {
		Code:      "final_negative",
		Title:     "Финальные отрицательные",
		BadgeType: "warning",
	},
}

var appealStatuses = map[string]AppealStatus{
	"new":               {Code: "new", Title: "Новое", Group: statusGroups["initial"]},
	"received":          {Code: "received", Title: "Поступило", Group: statusGroups["initial"]},
	"registered":        {Code: "registered", Title: "Зарегистрировано", Group: statusGroups["initial"]},
	"in_progress":       {Code: "in_progress", Title: "В работе", Group: statusGroups["working"]},
	"under_review":      {Code: "under_review", Title: "На рассмотрении", Group: statusGroups["working"]},
	"under_check":       {Code: "under_check", Title: "На проверке", Group: statusGroups["working"]},
	"extended":          {Code: "extended", Title: "Срок продлен", Group: statusGroups["working"]},
	"rework":            {Code: "rework", Title: "На доработке", Group: statusGroups["working"]},
	"forwarded":         {Code: "forwarded", Title: "Перенаправлено", Group: statusGroups["intermediate"]},
	"approval":          {Code: "approval", Title: "На согласовании", Group: statusGroups["intermediate"]},
	"signature":         {Code: "signature", Title: "На подписи", Group: statusGroups["intermediate"]},
	"awaiting_response": {Code: "awaiting_response", Title: "Ожидает ответа", Group: statusGroups["intermediate"]},
	"reviewed":          {Code: "reviewed", Title: "Рассмотрено", Group: statusGroups["final_positive"]},
	"answered":          {Code: "answered", Title: "Ответ направлен", Group: statusGroups["final_positive"]},
	"satisfied":         {Code: "satisfied", Title: "Удовлетворено", Group: statusGroups["final_positive"]},
	"approved":          {Code: "approved", Title: "Одобрено", Group: statusGroups["final_positive"]},
	"completed":         {Code: "completed", Title: "Исполнено", Group: statusGroups["final_positive"]},
	"rejected":          {Code: "rejected", Title: "Отказано", Group: statusGroups["final_negative"]},
}

var appealTransitions = map[string][]string{
	"new":               {"received", "registered", "in_progress", "rejected"},
	"received":          {"registered", "in_progress", "forwarded", "rejected"},
	"registered":        {"in_progress", "under_review", "forwarded", "rejected"},
	"in_progress":       {"under_review", "under_check", "extended", "rework", "reviewed", "forwarded", "rejected"},
	"under_review":      {"under_check", "approval", "awaiting_response", "reviewed", "rejected"},
	"under_check":       {"approval", "signature", "reviewed", "rework", "rejected"},
	"extended":          {"in_progress", "under_review", "reviewed", "rejected"},
	"rework":            {"in_progress", "under_review", "rejected"},
	"forwarded":         {"under_review", "awaiting_response", "reviewed", "rejected"},
	"approval":          {"signature", "approved", "rejected"},
	"signature":         {"answered", "approved", "rejected"},
	"awaiting_response": {"answered", "reviewed", "rejected"},
	"reviewed":          {"answered", "satisfied", "completed"},
	"answered":          {"satisfied", "completed"},
	"approved":          {"completed"},
	"satisfied":         {},
	"completed":         {},
	"rejected":          {},
}

func GetAppealStatus(code string) (AppealStatus, bool) {
	status, ok := appealStatuses[code]
	return status, ok
}

func ListAppealStatuses() []AppealStatus {
	items := make([]AppealStatus, 0, len(appealStatuses))
	for _, item := range appealStatuses {
		items = append(items, item)
	}
	return items
}

func AllowedTransitions(code string, isAdmin bool) []AppealStatus {
	current, ok := appealStatuses[code]
	if !ok {
		return nil
	}

	allowed := make([]AppealStatus, 0)
	for _, nextCode := range appealTransitions[code] {
		nextStatus, exists := appealStatuses[nextCode]
		if !exists {
			continue
		}
		allowed = append(allowed, nextStatus)
	}

	if isAdmin && (current.Group.Code == "final_positive" || current.Group.Code == "final_negative") {
		for _, nextCode := range []string{"in_progress", "under_review", "under_check", "rework"} {
			if nextStatus, exists := appealStatuses[nextCode]; exists {
				allowed = append(allowed, nextStatus)
			}
		}
	}

	return allowed
}

func IsAllowedTransition(from string, to string, isAdmin bool) bool {
	for _, status := range AllowedTransitions(from, isAdmin) {
		if status.Code == to {
			return true
		}
	}
	return false
}
