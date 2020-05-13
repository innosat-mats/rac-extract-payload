package innosat

// SourcePackageServiceType service type types
type SourcePackageServiceType uint8

const (
	// TelecommandVerification ...
	TelecommandVerification SourcePackageServiceType = 1
	// DeviceCommandDistribution ...
	DeviceCommandDistribution SourcePackageServiceType = 2
	// HousekeepingDiagnosticDataReporting ...
	HousekeepingDiagnosticDataReporting SourcePackageServiceType = 3
	// ParameterStatisticsReporting ...
	ParameterStatisticsReporting SourcePackageServiceType = 4
	// EventReporting ...
	EventReporting SourcePackageServiceType = 5
	// MemoryManagement ...
	MemoryManagement SourcePackageServiceType = 6
	// FunctionManagement ...
	FunctionManagement SourcePackageServiceType = 8
	// TimeManagement ...
	TimeManagement SourcePackageServiceType = 9
	// OnboardOperationsScheduling ...
	OnboardOperationsScheduling SourcePackageServiceType = 11
	// OnboardMonitoring ...
	OnboardMonitoring SourcePackageServiceType = 12
	// LargeDataTransfer ...
	LargeDataTransfer SourcePackageServiceType = 13
	// PacketForwardingControl ...
	PacketForwardingControl SourcePackageServiceType = 14
	// OnboardStorageandRetrieval ...
	OnboardStorageandRetrieval SourcePackageServiceType = 15
	// Test ...
	Test SourcePackageServiceType = 17
	// OnboardOperationsProcedure ...
	OnboardOperationsProcedure SourcePackageServiceType = 18
	// EventAction ...
	EventAction SourcePackageServiceType = 19
	// TransparentDataTransfer ...
	TransparentDataTransfer SourcePackageServiceType = 128
	// OnBoardParameterManagement ...
	OnBoardParameterManagement SourcePackageServiceType = 129
)

// SourcePackageServiceSubtype service subtype type (different meaning for each service)
type SourcePackageServiceSubtype uint8

const (
	// TCExecFailure subtype for Telecommand Execution Report - Failure
	TCExecFailure SourcePackageServiceSubtype = 8
	// TCAcceptSuccess subtype for Telecommand Failure Report - success
	TCAcceptSuccess SourcePackageServiceSubtype = 1
	// TCAcceptFailure subtype for Telecommand Failure Report - Failure
	TCAcceptFailure SourcePackageServiceSubtype = 2
	// TCExecSuccess subtype for Telecommand Execution Report - Success
	TCExecSuccess SourcePackageServiceSubtype = 7
)
