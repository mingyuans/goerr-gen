package code

// Job errors. module code = 01
const (
	// ErrSpannerUnknown - 500: DB Unknown error.
	ErrSpannerUnknown uint32 = iota + SystemCode*100*1000 + JobModuleCode*1000
	// ErrJobInvalid - 422: Job is invalid.
	ErrJobInvalid uint32 = iota + SystemCode*100*1000 + JobModuleCode*1000
	// ErrJobExists - 409: Job already exists, please check your unique key.
	ErrJobExists
	// ErrJobNotFound - 404: Job not found.
	ErrJobNotFound
	// ErrJobRedisUnknown - 500: DB Unknown error.
	ErrJobRedisUnknown
	// ErrJobCancelIsNotAllowed - 403: The current status is not allowed to be canceled.
	ErrJobCancelIsNotAllowed
	// ErrJobSQLBuilderUnknown - 500: SQL Builder Unknown error.
	ErrJobSQLBuilderUnknown
	// ErrLockFailed - 500: Lock Failed, please try again later.
	ErrLockFailed
	// ErrRetryingJobIsNotAllowed - 403: The job is not allowed to be retried.
	ErrRetryingJobIsNotAllowed
	// ErrRunningJobIsNotAllowed - 403: The job is not allowed to be running.
	ErrRunningJobIsNotAllowed
	// ErrJobIsAlreadyRunning - 403: The job is already running.
	ErrJobIsAlreadyRunning
)
