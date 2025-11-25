package keystore

const (
	// Health check action
	ActionPing = "ping"

	// Device key related actions
	ActionGetDeviceKey    = "getdevicekey"
	ActionSaveDeviceKey   = "savedevicekey"
	ActionDeleteDeviceKey = "deletedevicekey"

	// Session code related actions
	ActionGetSessionCode = "getsessioncode"

	// related to signup flow
	ActionSignAlias       = "signalias"
	ActionSaveSessionCode = "savesessioncode"

	// related to login flow
	ActionSignAliasWithTimestamp = "signaliaswithtimestamp"
	ActionSignChallengeToken     = "signchallengetoken"

	// related to login on another device
	ActionGenerateKeypair    = "generatekeypair"
	ActionGetPublicKey       = "getpublickey"
	ActionGetServerPublicKey = "getserverpubkey"
)
