package transport

type (
	Event struct {
		Origin EventOrigin
		Body   string

		Transport
	}

	EventOrigin struct {
		ID     string
		Sender EventOriginSender
	}

	EventOriginSender struct {
		ID   string
		Name string
	}
)
