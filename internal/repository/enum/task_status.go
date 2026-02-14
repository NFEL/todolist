package enum

type TaskStatus int

const (
	Created TaskStatus = iota
	Started
	Done
	Failed
	Delayed
	Canceled
)

func (t TaskStatus) String() string {
	switch t {
	case Created:
		return "Created"
	case Started:
		return "Started"
	case Done:
		return "Done"
	case Failed:
		return "Failed"
	case Delayed:
		return "Delayed"
	case Canceled:
		return "Canceled"
	default:
		return ""
	}
}
