package events

import "github.com/kentalee/eventbus"

const (
	PriorityPRE_ParseEnvArg  eventbus.Priority = 0
	PriorityPRE_InitLogger   eventbus.Priority = 10
	PriorityPRE_RegisterFlag eventbus.Priority = 20
)

const (
	PriorityINI_ParseFlag      eventbus.Priority = 0
	PriorityINI_CollectDefines eventbus.Priority = 5
	PriorityINI_CheckFlag      eventbus.Priority = 10
	PriorityINI_ParseConfig    eventbus.Priority = 20
)

const (
	PriorityAFT_CheckConfig eventbus.Priority = 10
)

const (
	PrioritySHU_AppShutdown eventbus.Priority = 64
	PrioritySHU_FlushLogger eventbus.Priority = 128
)
