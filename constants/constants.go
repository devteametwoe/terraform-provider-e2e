package constants

const WAIT_TIMEOUT = 3 // Waiting for change state in Seconds

const PREFIX_C2_NODE = "C2"

// node
//
//	i) status
var NODE_STATUS = map[string]string{
	"RUNNING":      "Running",
	"REINSTALLING": "Reinstalling",
	"CREATING":     "Creating",
	"FAILED":       "Failed",
	"POWERED_OFF":  "Powered off",
	"SAVING":       "Saving",
}

// ii) power_status
var NODE_POWER_STATUS = map[string]string{
	"ON":  "power_on",
	"OFF": "power_off",
}

// iii) lcm_state
var NODE_LCM_STATE = map[string]string{
	"HOTPLUG_PROLOG_POWEROFF": "HOTPLUG_PROLOG_POWEROFF",
	"HOTPLUG_EPILOG_POWEROFF": "HOTPLUG_EPILOG_POWEROFF",
	"HOTPLUG":                 "Hotplug",
}

// block storage

// i) action type
var BLOCK_STORAGE_ACTION = map[string]string{
	"ATTACH": "attach",
	"DETACH": "detach",
}

var BLOCK_STORAGE_STATUS = map[string]string{
	"ATTACHED":  "Attached",
	"SAVING":    "Saving",
	"CREATING":  "Creating",
	"AVAILABLE": "Available",
	"ERROR":     "ERROR",
}
