package constants

const WAIT_TIMEOUT = 3 // Waiting for change state in Seconds

const PREFIX_C2_NODE = "C2"

// node
//
//	i) status
const NODE_STATUS_RUNNING = "Running"
const NODE_STATUS_REINSTALLING = "Reinstalling"
const NODE_STATUS_CREATING = "Creating"
const NODE_STATUS_FAILED = "Failed"
const NODE_STATUS_POWERED_OFF = "Powered off"
const NODE_STATUS_SAVING = "Saving"

// ii) power_status
const NODE_POWER_STATUS_ON = "power_on"
const NODE_POWER_STATUS_OFF = "power_off"

// iii) lcm_state
const HOTPLUG_PROLOG_POWEROFF = "HOTPLUG_PROLOG_POWEROFF"
const HOTPLUG_EPILOG_POWEROFF = "HOTPLUG_EPILOG_POWEROFF"
const HOTPLUG = "Hotplug"

// block storage

// i) action type
const BLOCK_STORAGE_ACTION_ATTACH = "create"
const BLOCK_STORAGE_ACTION_DETACH = "detach"
